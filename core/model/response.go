package model

import (
	"time"
)

type Response struct {
	Headers    map[string]HeaderValues
	Body       []byte
	Status     string
	StatusCode int
	Proto      string
	Size       int64
	ReceivedAt time.Time
	Time       time.Duration
}
