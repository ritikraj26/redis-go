package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func rpushHandler(conn net.Conn, args []string) {
	if len(args) < 3 {
		resp.WriteError(conn, "Too few arguments for RPUSH")
		return
	}

	val := store.List[args[1]]
	for i := 2; i < len(args); i++ {
		val = append(val, args[i])
	}
	store.List[args[1]] = val

	resp.WriteInteger(conn, uint32(len(val)))
}
