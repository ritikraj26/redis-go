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

	key := args[1]
	val := store.List[args[1]]

	store.Mu.Lock()
	defer store.Mu.Unlock()

	// _, exists := store.BlockList[args[1]]
	// if exists {
	// 	var responseArray []string
	// 	responseArray = append(responseArray, args[1])
	// 	responseArray = append(responseArray, args[2])

	// 	resp.WriteArray(conn, responseArray)
	// 	return
	// }

	for i := 2; i < len(args); i++ {
		val = append(val, args[i])
	}
	store.List[key] = val

	if chans, exists := store.BlockingClients[key]; exists && len(chans) > 0 {
		clientChan := chans[0]
		store.BlockingClients[key] = store.BlockingClients[key][1:]

		go func() {
			client := <-clientChan
			resp.WriteArray(client, []string{key, args[2]})
		}()
	}

	resp.WriteInteger(conn, uint32(len(val)))
}
