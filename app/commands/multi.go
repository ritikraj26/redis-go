package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"

)

func multiHandler(conn net.Conn, args []string) {
	if len(args) != 1 {
		resp.WriteError(conn, "ERR: invalid number of arguments for MULTI")
		return
	}

	store.Mu.Lock()
	defer store.Mu.Unlock()

	store.TransactionQueue[conn] = []string{}
	resp.WriteSimpleString(conn, "OK")
}