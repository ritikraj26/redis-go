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

	key := args[1]

	store.Mu.Lock()
	val, exists := store.List[key]
	if !exists || len(val) == 0 {
		store.Mu.Unlock()
		resp.WriteNullString(conn)
		return
	}

	if len(args) == 2 {
		elem := val[0]
		store.List[key] = val[1:]
		store.Mu.Unlock()
		resp.WriteBulkString(conn, elem)
		return
	}

	count, err := strconv.Atoi(args[2])
	if err != nil {
		store.Mu.Unlock()
		resp.WriteError(conn, "ERR value is not an integer")
		return
	}

	if count <= 0 {
		store.Mu.Unlock()
		resp.WriteNullString(conn)
		return
	}

	count = min(count, len(val))

	var elements []string
	for i := 0; i < count; i++ {
		elements = append(elements, val[i])
	}
	val = val[count:]
	store.List[key] = val
	store.Mu.Unlock()

	resp.WriteArray(conn, elements)
}
