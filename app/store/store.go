package store

import (
	"time"
	"net"
	"sync"
)

type Data struct {
	Value  string
	Expiry *time.Time
}

var Store = make(map[string]Data)

var (
	BlockingClients = make(map[string][]chan net.Conn)
	Mu              sync.Mutex
)