package commands

import (
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func lrangeHandler(conn net.Conn, args []string) {
	if len(args) < 4 {
		resp.WriteError(conn, "Too few arguments")
		return
	}

	key := args[1]

	store.Mu.Lock()
	val, exists := store.List[key]
	store.Mu.Unlock()

	if !exists || len(val) == 0 {
		resp.WriteBulkStringArray(conn, []string{}) // Return empty array if list doesn't exist or is empty
		return
	}

	var respArray = []string{}

	leftIndex, err := strconv.Atoi(args[2])
	if err != nil {
		resp.WriteError(conn, "Index for lrange not correct")
	}
	rightIndex, err := strconv.Atoi(args[3])
	if err != nil {
		resp.WriteError(conn, "Index for lrange not correct")
	}

	if leftIndex < 0 {
		leftIndex += len(val)
	}

	if rightIndex < 0 {
		rightIndex += len(val)
	}

	if leftIndex < 0 {
		leftIndex = 0
	}

	if rightIndex >= len(val) {
		rightIndex = len(val) - 1
	}

	if leftIndex > rightIndex || leftIndex >= len(val) {
		resp.WriteBulkStringArray(conn, []string{})
		return
	}

	for i := leftIndex; i <= rightIndex; i++ {
		respArray = append(respArray, val[i])
	}

	resp.WriteBulkStringArray(conn, respArray)
}
