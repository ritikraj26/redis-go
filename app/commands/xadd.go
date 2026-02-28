package commands

import (
	"net"
	"strings"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func parseID(id string) (int64, int64, bool) {
	parts := strings.SplitN(id, "-", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}

	ms, err1 := strconv.ParseInt(parts[0], 10, 64)
	seq, err2 := strconv.ParseInt(parts[1], 10, 64)

	if err1 != nil || err2 != nil {
		return 0, 0, false
	}

	return ms, seq, true
}

func isValidId(streamName string, id string) (bool, string) {
	ms, seq, ok := parseID(id)
	if !ok {
		return false, "ERR Invalid stream ID specified as stream command argument"
	}

	// RULE #1: 0-0 is ALWAYS invalid
	if ms == 0 && seq == 0 {
		return false, "ERR The ID specified in XADD must be greater than 0-0"
	}

	store.Mu.Lock()
	defer store.Mu.Unlock()

	stream, exists := store.Stream[streamName]

	// Empty or non-existing stream → OK (since id != 0-0)
	if !exists || len(stream) == 0 {
		return true, ""
	}

	// Compare with last ID
	last := stream[len(stream)-1]
	lastMs, lastSeq, _ := parseID(last.Id)

	if ms > lastMs || (ms == lastMs && seq > lastSeq) {
		return true, ""
	}

	return false, "ERR The ID specified in XADD is equal or smaller than the target stream top item"
}

func xaddHandler(conn net.Conn, args []string) {
	if len(args) < 5 || (len(args)-3)%2 != 0 {
		resp.WriteError(conn, "Wrong number of arguments for XADD")
		return
	}

	streamName := args[1]
	id := args[2]

	ok, err := isValidId(streamName, id)
	if !ok {
		resp.WriteError(conn, err)
		return
	}

	data := store.StreamEntry{
		Id:     id,
		Fields: make(map[string]string),
	}

	for i := 3; i < len(args); i += 2 {
		data.Fields[args[i]]=args[i+1]
	}

	store.Mu.Lock()
	store.Stream[streamName] = append(store.Stream[streamName], data)
	store.Mu.Unlock()

	resp.WriteBulkString(conn, id)
	return
}