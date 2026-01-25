package commands

import (
	"net"
	"strconv"

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

	if len(args) == 2 {
		elem := val[0]
		store.List[args[1]] = val[1:]
		resp.WriteBulkString(conn, elem)
		return
	}

	count, err := strconv.Atoi(args[2])
	if err != nil {
		resp.WriteError(conn, "ERR value is not an integer")
		return
	}

	if count <= 0 {
		resp.WriteNullString(conn)
		return
	}

	count = min(count, len(val))

	var elements []string
	for i := 0; i < count; i++ {
		elements = append(elements, val[i])
	}
	val = val[count:]
	store.List[args[1]] = val

	resp.WriteArray(conn, elements)
}
