package net

import (
	"echo-go/internal/core"
	"fmt"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"

	"context"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type Transport struct {
	manager *core.Manager
	host    host.Host
}

func NewTransport(manager *core.Manager) *Transport {
	return &Transport{
		manager: manager,
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
	return nil
}
