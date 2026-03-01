package commands

import (
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func getHandler(conn net.Conn, args []string) {
	if len(args) < 2 {
		resp.WriteError(conn, "Get requires an argument (key)")
		return
	}

	key := args[1]

	store.Mu.Lock()
	val, exists := store.Store[key]
	if !exists {
		store.Mu.Unlock()
		resp.WriteNullBulkString(conn)
		return
	}

	if val.Expiry != nil && time.Now().After(*val.Expiry) {
		delete(store.Store, key)
		store.Mu.Unlock()
		resp.WriteNullBulkString(conn)
		return
	}

	store.Mu.Unlock()
	resp.WriteBulkString(conn, val.Value)
}
