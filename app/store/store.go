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

var (
	BlockingClients = make(map[string][]chan string)
	Mu              sync.Mutex
)