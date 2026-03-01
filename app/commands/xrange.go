package commands

import (
	"fmt"
	"net"
	"strings"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func parseRangeId(id string, isStart bool) (int64, int64) {
	// If sequence number missing
	if !strings.Contains(id, "-") {
		ms, _ := strconv.ParseInt(id, 10, 64)
		if isStart {
			return ms, 0
		}
		return ms, int64(^uint64(0) >> 1) // max int64
	}

	parts := strings.SplitN(id, "-", 2)
	ms, _ := strconv.ParseInt(parts[0], 10, 64)
	seq, _ := strconv.ParseInt(parts[1], 10, 64)
	return ms, seq
}

func xrangeHandler(conn net.Conn, args []string) {
	if len(args) != 4 {
		resp.WriteError(conn, "ERR wrong number of arguments for XRANGE")
		return
	}

	key := args[1]
	start := args[2]
	end := args[3]

	startMs, startSeq := parseRangeId(start, true)
	endMs, endSeq := parseRangeId(end, false)

	store.Mu.Lock()
	stream := store.Stream[key]
	store.Mu.Unlock()

	// Return empty array if stream missing
	if len(stream) == 0 {
		resp.WriteArray(conn, []string{})
		return
	}

	// Collect matching entries
	results := []store.StreamEntry{}
	for _, entry := range stream {
		ms, seq := parseRangeId(entry.Id, true)

		if (ms > startMs || (ms == startMs && seq >= startSeq)) &&
		   (ms < endMs   || (ms == endMs   && seq <= endSeq)) {
			results = append(results, entry)
		}
	}

	// Outer array
	fmt.Fprintf(conn, "*%d\r\n", len(results))

	for _, entry := range results {
		// Each entry: [id, [fields]]
		fmt.Fprintf(conn, "*2\r\n")
		resp.WriteBulkString(conn, entry.Id)

		fmt.Fprintf(conn, "*%d\r\n", len(entry.Fields))
		for _, f := range entry.Fields {
			resp.WriteBulkString(conn, f)
		}
	}
}