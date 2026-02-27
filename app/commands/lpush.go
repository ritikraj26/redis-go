package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func lpushHandler(conn net.Conn, args []string) {
	if len(args) < 3 {
		resp.WriteError(conn, "Too few arguments for LPUSH")
		return
	}

	key := args[1]

	val := store.List[key]
	for i := 2; i < len(args); i++ {
		val = append([]string{args[i]}, val...)
	}
	store.List[key] = val

	resp.WriteInteger(conn, uint32(len(val)))
}
