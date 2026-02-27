package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func llenHandler(conn net.Conn, args []string) {
	if len(args) < 2 {
		resp.WriteError(conn, "Too few arguments")
		return
	}

	key := args[1]

	val, exists := store.List[key]
	if !exists || len(val) == 0 {
		resp.WriteInteger(conn, 0)
		return
	}

	resp.WriteInteger(conn, uint32(len(val)))
}
