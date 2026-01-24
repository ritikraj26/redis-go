package store

import "time"

type Data struct {
	Value  string
	Expiry *time.Time
}

var Store = make(map[string]Data)
