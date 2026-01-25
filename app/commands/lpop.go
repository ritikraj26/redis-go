package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func lpopHandler(conn net.Conn, args []string) {
	if len(args) < 2 {
		resp.WriteError(conn, "Too few arguments for LPOP")
		return
	}

	val, exists := store.List[args[1]]
	if !exists || len(val) == 0 {
		resp.WriteNullString(conn)
		return
	}

	element := val[0]
	val = val[1:]
	store.List[args[1]] = val
	resp.WriteBulkString(conn, element)
}
