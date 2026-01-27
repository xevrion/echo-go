package main

import (
	"echo-go/internal/core"
	"echo-go/internal/net"
	"echo-go/internal/ui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config := core.Config{
		Username: "xevrion",
		Port:     8080,
	}

	manager := core.NewManager(config)
	transport := net.NewTransport(manager)
	defer transport.Stop() // ensure transport stops on exit
	transport.Start()

	model := tui.NewModel(manager)
	program := tea.NewProgram(model)

	// bridge: core → bubbletea
	go func() {
		for event := range manager.Events() {
			program.Send(event)
		}
	}()

	// test message
	go func() {
		manager.Send("hello from echo-go")
	}()

	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
