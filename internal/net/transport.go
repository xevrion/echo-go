package net

import (
	"echo-go/internal/core"

	"github.com/libp2p/go-libp2p/core/host"
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
	return nil
}

func (t *Transport) Stop() error {
	return nil
}
