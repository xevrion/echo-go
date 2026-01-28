package main

import (
	"echo-go/internal/core"
	"echo-go/internal/net"
	"echo-go/internal/ui/tui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	port := 8081
	if p := os.Getenv("PORT"); p != "" {
		fmt.Sscanf(p, "%d", &port)
	}

	config := core.Config{
		Username: "xevrion",
		Port:     port,
	}

	manager := core.NewManager(config)
	transport := net.NewTransport(manager)
	defer transport.Stop() // ensure transport stops on exit
	transport.Start()

	if len(os.Args) > 1 {
		transport.Connect(os.Args[1])
	}

	model := tui.NewModel(manager, transport)
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
