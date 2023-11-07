package event

import (
	"time"
)

type Message interface {
	Key() string
}

type GiggleCreatedMessage struct {
	ID        string
	Body      string
	CreatedAt time.Time
}

func (m *GiggleCreatedMessage) Key() string {
	return "giggle.created"
}
