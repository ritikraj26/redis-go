package commands

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func parseId(id string) (int64, int64, bool) {
	parts := strings.SplitN(id, "-", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}
	ms, e1 := strconv.ParseInt(parts[0], 10, 64)
	seq, e2 := strconv.ParseInt(parts[1], 10, 64)
	return ms, seq, e1 == nil && e2 == nil
}

func generateId(streamName, id string) (string, string) {
	store.Mu.Lock()
	defer store.Mu.Unlock()

	stream := store.Streams[streamName]

	// *
	if id == "*" {
		ms := time.Now().UnixMilli()
		seq := int64(0)
		if len(stream) > 0 {
			lms, lseq, _ := parseId(stream[len(stream)-1].Id)
			if lms == ms {
				seq = lseq + 1
			}
		}
		return fmt.Sprintf("%d-%d", ms, seq), ""
	}

	// <ms>-*
	if strings.HasSuffix(id, "-*") {
		msPart := strings.TrimSuffix(id, "-*")
		ms, err := strconv.ParseInt(msPart, 10, 64)
		if err != nil {
			return "", "ERR Invalid stream ID specified as stream command argument"
		}
		seq := int64(0)
		if len(stream) > 0 {
			lms, lseq, _ := parseId(stream[len(stream)-1].Id)
			if lms == ms {
				seq = lseq + 1
			}
		} else if ms == 0 {
			seq = 1
		}
		return fmt.Sprintf("%d-%d", ms, seq), ""
	}

	// explicit ID
	ms, seq, ok := parseId(id)
	if !ok {
		return "", "ERR Invalid stream ID specified as stream command argument"
	}

	if ms == 0 && seq == 0 {
		return "", "ERR The ID specified in XADD must be greater than 0-0"
	}
	if len(stream) > 0 {
		lms, lseq, _ := parseId(stream[len(stream)-1].Id)
		if ms < lms || (ms == lms && seq <= lseq) {
			return "", "ERR The ID specified in XADD is equal or smaller than the target stream top item"
		}
	}
	return id, ""
}

func xaddHandler(conn net.Conn, args []string) {
	if len(args) < 5 || (len(args)-3)%2 != 0 {
		resp.WriteError(conn, "ERR: wrong number of arguments for XADD")
		return
	}

	stream := args[1]
	id := args[2]

	newId, err := generateId(stream, id)
	if err != "" {
		resp.WriteError(conn, err)
		return
	}

	entry := store.StreamEntry{Id: newId}
	for i := 3; i < len(args); i += 2 {
		entry.Fields = append(entry.Fields, args[i], args[i+1])
	}

	store.Mu.Lock()
	store.Streams[stream] = append(store.Streams[stream], entry)

	// wake blocked XREAD clients
	waiters := store.StreamBlockingClients[stream]
	delete(store.StreamBlockingClients, stream)
	store.Mu.Unlock()

	for _, ch := range waiters {
		select {
		case ch <- struct{}{}:
		default:
		}
	}

	resp.WriteBulkString(conn, newId)
}