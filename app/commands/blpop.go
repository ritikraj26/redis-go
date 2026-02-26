package commands

import (
	"net"
	"time"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func blpopHandler(conn net.Conn, args []string) {
	if len(args) < 3 {
		resp.WriteError(conn, "Too few arguments")
		return
	}

	key := args[1]
	timeoutSeconds, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		resp.WriteError(conn, "Invalid timeout value")
	}

	store.Mu.Lock()
	// check if the list has elements
	if len(store.List[key]) > 0 {
		element := store.List[key][0]
		store.List[key] = store.List[key][1:] // removing the first element

		resp.WriteArray(conn, []string{key, element})
		return
	}

	// block the client if the list is empty
	clientChan := make(chan string, 1)
	store.BlockingClients[key] = append(store.BlockingClients[key], clientChan)
	store.Mu.Unlock()

	if timeoutSeconds == 0 {
		element := <-clientChan
		resp.WriteArray(conn, []string{key, element})
		return
	}

	select {
	case element := <-clientChan:
		resp.WriteArray(conn, []string{key, element})
	case <-time.After(time.Duration(timeoutSeconds * float64(time.Second))):
		store.Mu.Lock()
		waiters := store.BlockingClients[key]
		for i, c := range waiters {
			if c == clientChan {
				store.BlockingClients[key] = append(waiters[:i], waiters[i+1:]...)
				break
			}
		}
		store.Mu.Unlock()

		// RESP null array
		resp.WriteNullArray(conn)
	}
}
