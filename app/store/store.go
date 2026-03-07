package store

import (
	"sync"
	"time"
	"net"
)

type Data struct {
	Value  string
	Expiry *time.Time
}

type StreamEntry struct {
	Id     string
	Fields []string
}

var (
	Store   = make(map[string]Data)
	List    = make(map[string][]string)
	Streams = make(map[string][]StreamEntry)

	// BLPOP uses this (string payload)
	ListBlockingClients = make(map[string][]chan string)

	// XREAD BLOCK uses this (signal only)
	StreamBlockingClients = make(map[string][]chan struct{})

	// MULTI uses this
	TransactionQueue = make(map[net.Conn][]string)

	Mu sync.Mutex
)