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
	if len(args) != 4 {
		resp.WriteError(conn, "ERR wrong number of arguments for XREAD")
		return
	}

	key := args[2]
	startId := args[3]

	startMs, startSeq, ok := parseReadId(startId)
	if !ok {
		resp.WriteError(conn, "ERR Invalid stream ID specified as stream command argument")
		return
	}

	store.Mu.Lock()
	stream, exists := store.Stream[key]
	store.Mu.Unlock()

	if !exists || len(stream) == 0 {
		resp.WriteBulkStringArray(conn, []string{})
		return
	}

	var entries []store.StreamEntry
	for _, entry := range stream {
		ms, seq, _ := parseReadId(entry.Id)
		if ms > startMs || (ms == startMs && seq >= startSeq) {
			entries = append(entries, entry)
		}
	}

	if len(entries) == 0 {
		resp.WriteBulkStringArray(conn, []string{})
		return
	}

	// single stream
	fmt.Fprintf(conn, "*1\r\n")
	// two elements per stream
	fmt.Fprintf(conn, "*2\r\n")

	resp.WriteBulkString(conn, key)

	fmt.Fprintf(conn, "*%d\r\n", len(entries))
	for _, entry := range entries {
		fmt.Fprintf(conn, "*2\r\n")
		resp.WriteBulkString(conn, entry.Id)

		fmt.Fprintf(conn, "*%d\r\n", len(entry.Fields))
		for _, f := range entry.Fields {
			resp.WriteBulkString(conn, f)
		}
	}
}