package store

import (
	"sync"
	"time"
)

type Data struct {
	Value  string
	Expiry *time.Time
}

type StreamEntry struct {
	Id     string
	Fields map[string]string
}

var (
	Store  = make(map[string]Data)
	List   = make(map[string][]string)
	Stream = make(map[string][]StreamEntry)

	BlockingClients = make(map[string][]chan string)

	Mu sync.Mutex
)