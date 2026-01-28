package tui

import (
	"echo-go/internal/core"
	"echo-go/internal/net"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	manager   *core.Manager
	messages  []core.Message
	ready     bool
	input     string
	transport *net.Transport
	peers     []core.Peer
	cursor    int
}

func NewModel(manager *core.Manager, transport *net.Transport) *Model {
	return &Model{
		manager:   manager,
		messages:  []core.Message{},
		ready:     false,
		transport: transport,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case core.Event:
		if msg.Type == "message" {
			if message, ok := msg.Payload.(core.Message); ok {
				m.messages = append(m.messages, message)
			}
		}
		if msg.Type == core.EventPeerFound {
			m.peers = append(m.peers, msg.Payload.(core.Peer))
		}

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEnter:
			if m.input != "" {
				// normal message send
				m.transport.Send(m.input)
				m.input = ""
			} else if len(m.peers) > 0 {
				// connect to selected peer
				selected := m.peers[m.cursor]
				m.transport.Connect(selected.ID)
			}

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}

		case tea.KeyDown:
			if m.cursor < len(m.peers)-1 {
				m.cursor++
			}

		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m *Model) View() string {
	s := "Online:\n"
	for i, p := range m.peers {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		s += prefix + p.Name + "\n"
	}

	s += "\n"
	for _, msg := range m.messages {
		s += msg.Sender + ": " + msg.Text + "\n"
	}
	s += "\n> " + m.input
	return s
}
