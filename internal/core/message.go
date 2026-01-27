package core

import (
	"time"
)

type Message struct {
	Sender string
	Text   string
	Time   time.Time
}
