package commands

import (
	"fmt"
	"net"
	"time"
	"strings"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func parseId(id string) (int64, int64, bool) {
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
	ms, seq, ok := parseId(id)
	if !ok {
		return false, "ERR Invalid stream ID specified as stream command argument"
	}

	if ms == 0 && seq == 0 {
		return false, "ERR The ID specified in XADD must be greater than 0-0"
	}

	stream := store.Streams[streamName]
	if len(stream) == 0 {
		return true, ""
	}

	lastMs, lastSeq, _ := parseId(stream[len(stream)-1].Id)
	if ms > lastMs || (ms == lastMs && seq > lastSeq) {
		return true, ""
	}

	return false, "ERR The ID specified in XADD is equal or smaller than the target stream top item"
}

func generateId(streamName string, id string) (string, string) {
	store.Mu.Lock()
	defer store.Mu.Unlock()

	stream := store.Streams[streamName]

	// Case 1: *
	if id == "*" {
		ms := time.Now().UnixMilli()
		seq := int64(0)

		if len(stream) > 0 {
			lastMs, lastSeq, _ := parseId(stream[len(stream)-1].Id)
			if lastMs == ms {
				seq = lastSeq + 1
			}
		}
		return fmt.Sprintf("%d-%d", ms, seq), ""
	}

	// Case 2: <ms>-*
	if strings.HasSuffix(id, "-*") {
		msPart := strings.TrimSuffix(id, "-*")
		ms, err := strconv.ParseInt(msPart, 10, 64)
		if err != nil {
			return "", "ERR Invalid stream ID specified as stream command argument"
		}

		var seq int64

		if len(stream) == 0 {
			if ms == 0 {
				seq = 1
			} else {
				seq = 0
			}
		} else {
			lastMs, lastSeq, _ := parseId(stream[len(stream)-1].Id)
			if lastMs == ms {
				seq = lastSeq + 1
			} else {
				seq = 0
			}
		}

		return fmt.Sprintf("%d-%d", ms, seq), ""
	}

	// Case 3: explicit ID
	ok, err := isValidId(streamName, id)
	if !ok {
		return "", err
	}

	return id, ""
}

func xaddHandler(conn net.Conn, args []string) {
	if len(args) < 5 || (len(args)-3)%2 != 0 {
		resp.WriteError(conn, "Wrong number of arguments for XADD")
		return
	}

	streamName := args[1]
	id := args[2]

	newId, err := generateId(streamName, id)
	if err != "" {
		resp.WriteError(conn, err)
		return
	}

	data := store.StreamEntry{
		Id:     newId,
		Fields: []string{},
	}

	// Append fields in order: key, value, key, value...
	for i := 3; i < len(args); i += 2 {
		data.Fields = append(data.Fields, args[i], args[i+1])
	}

	store.Mu.Lock()
	store.Streams[streamName] = append(store.Streams[streamName], data)
	store.Mu.Unlock()

	resp.WriteBulkString(conn, newId)
}