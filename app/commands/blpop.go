package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func blpopHandler(conn net.Conn, args []string) {
	if len(args) < 3 {
		resp.WriteError(conn, "Too few arguments")
		return
	}

	key := args[1]

	store.Mu.Lock()
	defer store.Mu.Unlock()

	// check if the list has elements
	if len(store.List[key]) > 0 {
		element := store.List[key][0]
		store.List[key] = store.List[key][1:] // removing the first element

		resp.WriteArray(conn, []string{key, element})
		return
	}

	// block the client if the list is empty, now indefinitely
	clientChan := make(chan net.Conn, 1)
	store.BlockingClients[key] = append(store.BlockingClients[key], clientChan)

	//unlock the mutex while waiting
	store.Mu.Unlock()
	clientChan <- conn
	store.Mu.Lock()
}
