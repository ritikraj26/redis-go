package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func xaddHandler(conn net.Conn, args []string) {
	if len(args) < 5 || (len(args)-3)%2 != 0 {
		resp.WriteError(conn, "Wrong number of arguments for XADD")
		return
	}

	streamName := args[1]
	id := args[2]

	data := store.StreamEntry{
		Id:     id,
		Fields: make(map[string]string),
	}

	for i := 3; i < len(args); i += 2 {
		data.Fields[args[i]]=args[i+1]
	}

	store.Stream[streamName] = append(store.Stream[streamName], data)

	resp.WriteBulkString(conn, id)
	return
}