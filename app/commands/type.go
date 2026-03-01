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

	store.Mu.Lock()

	_, isString := store.Store[key]
	_, isList := store.List[key]
	_, isStream := store.Streams[key]

	store.Mu.Unlock()

	if isString {
		resp.WriteSimpleString(conn, "string")
		return
	}

	if isList {
		resp.WriteSimpleString(conn, "list")
		return
	}

	if isStream {
		resp.WriteSimpleString(conn, "stream")
		return
	}

	resp.WriteSimpleString(conn, "none")
}