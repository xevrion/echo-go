package main

import (
	"fmt"
	"time"

	"echo-go/internal/core"
)

func main() {
	config := core.Config{
		Username: "xevrion",
		Port:     8080,
	}

	manager := core.NewManager(config)

	// listen for events (runs in background)
	go func() {
		for event := range manager.Events() {
			fmt.Println(event.Type, event.Payload)
		}
	}()

	// test actions
	manager.Send("Hello, World!")

	// prevent program from exiting
	time.Sleep(1 * time.Second)
}
