package commands

import (
	"fmt"
	"net"
	"strings"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func parseReadId(id string) (int64, int64, bool) {
	if !strings.Contains(id, "-") {
		return 0, 0, false
	}

	parts := strings.SplitN(id, "-", 2)
	ms, err1 := strconv.ParseInt(parts[0], 10, 64)
	seq, err2 := strconv.ParseInt(parts[1], 10, 64)

	if err1 != nil || err2 != nil {
		return 0, 0, false
	}

	return ms, seq, true
}

func xreadHandler(conn net.Conn, args []string) {
	if len(args) < 4 {
		resp.WriteError(conn, "ERR wrong number of arguments for XREAD")
		return
	}

	type StreamData struct {
		Key     string
		Entries []store.StreamEntry
	}

	streamCount := (len(args) - 2) / 2
	keys := args[2 : 2 + streamCount]
	ids := args[2 + streamCount : ]

	var result []StreamData

	for i := 0; i < streamCount; i++ {
		key := keys[i]
		startId := ids[i]

		startMs, startSeq, ok := parseReadId(startId)
		if !ok {
			resp.WriteError(conn, "ERR Invalid stream ID specified as stream command argument")
			return
		}

		store.Mu.Lock()
		stream, exists := store.Streams[key]
		store.Mu.Unlock()

		if !exists || len(stream) == 0 {
			continue
		}

		var entries []store.StreamEntry
		for _, entry := range stream {
			ms, seq, _ := parseReadId(entry.Id)
			if ms > startMs || (ms == startMs && seq >= startSeq) {
                entries = append(entries, entry)
            }
		}

		if len(entries) > 0 {
			result = append(result, StreamData{Key: key, Entries: entries})
		}
	}

	if len(result) == 0 {
        resp.WriteBulkStringArray(conn, []string{})
        return
    }

	// multiple stream
	fmt.Fprintf(conn, "*%d\r\n", len(result))
	for _, entries := range result {
		fmt.Fprintf(conn, "*2\r\n")
		resp.WriteBulkString(conn, entries.Key)
		fmt.Fprintf(conn, "*%d\r\n", len(entries.Entries))
		for _, entry := range entries.Entries {
			fmt.Fprintf(conn, "*2\r\n")
			resp.WriteBulkString(conn, entry.Id)

			fmt.Fprintf(conn, "*%d\r\n", len(entry.Fields))
			for _, f := range entry.Fields {
				resp.WriteBulkString(conn, f)
			}
		}
	}
}