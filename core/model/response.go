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
	TraceInfo  TraceInfo
}

type TraceInfo struct {
	DNSLookup      time.Duration
	ConnTime       time.Duration
	TCPConnTime    time.Duration
	TLSHandshake   time.Duration
	ServerTime     time.Duration
	ResponseTime   time.Duration
	TotalTime      time.Duration
	IsConnReused   bool
	IsConnWasIdle  bool
	ConnIdleTime   time.Duration
	RequestAttempt int
	RemoteAddr     string
}
