package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func echoHandler(conn net.Conn, args []string) {
	if len(args) < 2 {
		resp.WriteError(conn, "ECHO requires an argument")
		return
	}

	key := args[1]
	resp.WriteBulkString(conn, key)
}
