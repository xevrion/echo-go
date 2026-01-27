package core

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
