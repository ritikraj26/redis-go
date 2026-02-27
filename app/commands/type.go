package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func typeHandler(conn net.Conn, args []string) {
	if len(args) < 2 {
		resp.WriteError(conn, "Too few arguments for TYPE")
		return;
	}

	key := args[1]

	_, exists := store.Store[key]
	if !exists {
		resp.WriteSimpleString(conn, "none")
		return
	}

	resp.WriteSimpleString(conn, "string")
	return
}