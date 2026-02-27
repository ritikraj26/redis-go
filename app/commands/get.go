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

	val, exists := store.Store[key]
	if !exists {
		resp.WriteNullString(conn)
		return
	}

	if val.Expiry != nil && time.Now().After(*val.Expiry) {
		delete(store.Store, key)
		resp.WriteNullString(conn)
		return
	}

	resp.WriteBulkString(conn, val.Value)
}
