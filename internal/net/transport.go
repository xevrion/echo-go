package net

import "echo-go/internal/core"

type Transport struct {
	manager *core.Manager
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
