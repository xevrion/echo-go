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

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEnter:
			if m.input != "" {
				m.manager.Send(m.input)
				m.transport.Send(m.input)
				m.input = ""
			}

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m *Model) View() string {
	s := ""
	for _, msg := range m.messages {
		s += msg.Sender + ": " + msg.Text + "\n"
	}
	s += "\n> " + m.input
	return s
}
