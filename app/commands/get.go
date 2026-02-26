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

	val, exists := store.Store[args[1]]
	if !exists {
		resp.WriteNullString(conn)
		return
	}

	if val.Expiry != nil && time.Now().After(*val.Expiry) {
		delete(store.Store, args[1])
		resp.WriteNullString(conn)
		return
	}

	resp.WriteBulkString(conn, val.Value)
}
