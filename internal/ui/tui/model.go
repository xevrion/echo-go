package tui

import (
	"echo-go/internal/core"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	manager  *core.Manager
	messages []core.Message
	ready    bool
}

func NewModel(manager *core.Manager) *Model {
	return &Model{
		manager:  manager,
		messages: []core.Message{},
		ready:    false,
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
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) View() string {
	s := "Messages:\n"
	for _, msg := range m.messages {
		s += msg.Sender + ": " + msg.Text + "\n"
	}
	s += "\nPress q to quit.\n"
	return s
}
