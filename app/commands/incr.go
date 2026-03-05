package commands

import (
	"net"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func incrHandler(conn net.Conn, args []string) {
	if len(args) != 2 {
		resp.WriteError(conn, "ERR: wrong number of arguments for INCR")
		return
	}

	key := args[1]

	store.Mu.Lock()
	defer store.Mu.Unlock()

	data, exists := store.Store[key]

	var intVal int
	var err error

	if !exists {
		intVal = 1
	} else {
		intVal, err = strconv.Atoi(data.Value)
		if err != nil {
			resp.WriteError(conn, "ERR value is not an integer or out of range")
			return
		}
		intVal++
	}
	data.Value = strconv.Itoa(intVal)
	store.Store[key] = data

	resp.WriteInteger(conn, uint32(intVal))
}