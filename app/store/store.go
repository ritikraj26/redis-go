package store

import (
	"time"
	"sync"
)

type Data struct {
	Value  string
	Expiry *time.Time
}

var Store = make(map[string]Data)

var List = make(map[string][]string)

var (
	BlockingClients = make(map[string][]chan string)
	Mu              sync.Mutex
)

type StreamEntry struct {
	Id     string
	Fields map[string]string
}
var Stream = make(map[string][]StreamEntry)