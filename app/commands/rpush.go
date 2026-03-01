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

	store.Mu.Lock()
	// defer store.Mu.Unlock()

	// If client is blocked
	if chans, exists := store.BlockingClients[key]; exists && len(chans) > 0 {
		clientChan := chans[0]
		store.BlockingClients[key] = store.BlockingClients[key][1:]
		store.Mu.Unlock()

		clientChan <- args[2]
		resp.WriteInteger(conn, 1)
		return
		// go func() {
		// 	client := <-clientChan
		// 	resp.WriteBulkStringArray(client, []string{key, args[2]})
		// }()
	}

	// No blocked client
	val := store.List[key]
	for i := 2; i < len(args); i++ {
		val = append(val, args[i])
	}
	store.List[key] = val

	store.Mu.Unlock()

	resp.WriteInteger(conn, uint32(len(val)))
}
