package main

import (
	"echo-go/internal/core"
	"echo-go/internal/net"
	"echo-go/internal/ui/tui"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	port := 8082
	if p := os.Getenv("PORT"); p != "" {
		fmt.Sscanf(p, "%d", &port)
	}

	name := os.Getenv("NAME")
	if name == "" {
		name = "user"
	}

	config := core.Config{
		Username: name,
		Port:     port,
	}

	manager := core.NewManager(config)
	transport := net.NewTransport(manager)
	defer transport.Stop()
	transport.Start()

	if len(os.Args) > 1 {
		time.Sleep(100 * time.Millisecond)
		transport.Connect(os.Args[1])
	}

	model := tui.NewModel(manager, transport)
	program := tea.NewProgram(model)

	go func() {
		for event := range manager.Events() {
			program.Send(event)
		}
	}()

	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
