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

    if _, exists := store.Store[key]; exists {
        resp.WriteSimpleString(conn, "string")
        return
    }

    if _, exists := store.List[key]; exists {
        resp.WriteSimpleString(conn, "list")
        return
    }

    if _, exists := store.Stream[key]; exists {
        resp.WriteSimpleString(conn, "stream")
        return
    }

	resp.WriteSimpleString(conn, "none")
	return
}