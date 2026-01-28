package net

import (
	"bufio"
	"context"
	"echo-go/internal/core"
	"fmt"
	"strings"
	"sync"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/multiformats/go-multiaddr"
)

const ProtocolID = "/echo-go/chat/1.0.0"

type Transport struct {
	manager *core.Manager
	host    host.Host
	streams map[peer.ID]network.Stream
	mu      sync.Mutex
}

func NewTransport(manager *core.Manager) *Transport {
	return &Transport{
		manager: manager,
		streams: make(map[peer.ID]network.Stream),
	}
}

func (t *Transport) Start() error {
	port := t.manager.Port()

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port),
		),
	)
	if err != nil {
		return err
	}

	t.host = h
	notifee := &discoveryNotifee{transport: t}

	service := mdns.NewMdnsService(h, "echo-go-mdns", notifee)
	if err := service.Start(); err != nil {
		return err
	}

	t.host.SetStreamHandler(ProtocolID, t.handleStream)

	fmt.Println("Peer ID:", h.ID())
	for _, addr := range h.Addrs() {
		fmt.Println("Listening on:", addr)
	}

	return nil
}

func (t *Transport) Stop() error {
	return nil
}

func (t *Transport) Connect(addr string) error {
	maddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return err
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}

	if err := t.host.Connect(context.Background(), *info); err != nil {
		return err
	}

	fmt.Println("Connected to:", info.ID)

	t.mu.Lock()
	_, exists := t.streams[info.ID]
	t.mu.Unlock()

	if exists {
		return nil
	}

	stream, err := t.host.NewStream(context.Background(), info.ID, ProtocolID)
	if err != nil {
		return err
	}

	t.mu.Lock()
	t.streams[info.ID] = stream
	t.mu.Unlock()

	go t.readStream(stream)

	return nil
}

func (t *Transport) handleStream(s network.Stream) {
	peerID := s.Conn().RemotePeer()

	t.mu.Lock()
	_, exists := t.streams[peerID]
	if exists {
		t.mu.Unlock()
		s.Close()
		return
	}
	t.streams[peerID] = s
	t.mu.Unlock()

	t.readStream(s)
}

func (t *Transport) readStream(s network.Stream) {
	peerID := s.Conn().RemotePeer()
	scanner := bufio.NewScanner(s)

	for scanner.Scan() {
		raw := scanner.Text()

		parts := strings.SplitN(raw, "|", 2)
		sender := peerID.String()
		msgText := raw

		if len(parts) == 2 {
			sender = parts[0]
			msgText = parts[1]
		}

		msg := core.Message{
			Sender: sender,
			Text:   msgText,
			Time:   time.Now(),
		}

		t.manager.Receive(msg)
	}

	t.mu.Lock()
	delete(t.streams, peerID)
	t.mu.Unlock()

	s.Close()
}

func (t *Transport) Send(text string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for peerID, stream := range t.streams {
		_, err := fmt.Fprintf(stream, "%s|%s\n", t.manager.Username(), text)
		if err != nil {
			delete(t.streams, peerID)
			stream.Close()
		}
	}
}

type discoveryNotifee struct {
	transport *Transport
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	addr := pi.Addrs[0].Encapsulate(
		multiaddr.StringCast(fmt.Sprintf("/p2p/%s", pi.ID)),
	)

	n.transport.manager.AddDiscoveredPeer(core.Peer{
		ID:   pi.ID.String(),
		Name: pi.ID.String(), // temporary name for now
	})

	n.transport.Connect(addr.String())
}

func (t *Transport) SendTo(peerID, text string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	stream, ok := t.streams[peer.ID(peerID)]
	if !ok {
		return
	}

	fmt.Fprintf(stream, "%s|%s\n", t.manager.Username(), text)
}
