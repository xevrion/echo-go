package net

import (
	"bufio"
	"echo-go/internal/core"
	"fmt"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"

	"context"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const ProtocolID = "/echo-go/chat/1.0.0"

type Transport struct {
	manager *core.Manager
	host    host.Host
	streams map[peer.ID]network.Stream
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
	fmt.Println("Dialing:", addr)

	fmt.Println("Connected to:", info.ID)
	return nil
}

func (t *Transport) handleStream(s network.Stream) {
	defer s.Close()

	scanner := bufio.NewScanner(s)

	for scanner.Scan() {
		text := scanner.Text()

		msg := core.Message{
			Sender: s.Conn().RemotePeer().String(),
			Text:   text,
			Time:   time.Now(),
		}

		t.manager.Receive(msg)
	}
}
