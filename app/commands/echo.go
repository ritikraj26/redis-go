package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func Echo(conn net.Conn, args []string) {
	if len(args) < 2 {
		resp.WriteError(conn, "ECHO requires an argument")
		return
	}
	resp.WriteBulkString(conn, args[1])
}
