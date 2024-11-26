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
	Options    map[string]interface{}
}

type TraceInfo struct {
	IsConnReused   bool
	IsConnWasIdle  bool
	DNSLookup      time.Duration
	ConnTime       time.Duration
	TCPConnTime    time.Duration
	TLSHandshake   time.Duration
	ServerTime     time.Duration
	ResponseTime   time.Duration
	TotalTime      time.Duration
	ConnIdleTime   time.Duration
	RequestAttempt int
	RemoteAddr     string
}

func (r *Response) HeadersAsMapString() map[string][]string {
	headers := make(map[string][]string)
	for k, v := range r.Headers {
		headers[k] = v
	}
	return headers
}
