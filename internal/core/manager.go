package core

import "time"

type Manager struct {
	config Config
	peers  map[string]Peer
	events chan Event
}

func NewManager(config Config) *Manager {
	return &Manager{
		config: config,
		peers:  make(map[string]Peer),
		events: make(chan Event, 100),
	}
}

func (manager *Manager) Events() chan Event {
	return manager.events
}

func (manager *Manager) Send(text string) {
	// Implementation for sending a message
	msg := Message{
		Text:   text,
		Sender: manager.config.Username,
		Time:   time.Now(),
	}

	manager.events <- Event{
		Type:    "message",
		Payload: msg,
	}

}

func (manager *Manager) Receive(msg Message) {
	manager.events <- Event{
		Type:    "message",
		Payload: msg,
	}
}
