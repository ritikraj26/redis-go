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

type StreamData struct {
	Key     string
	Entries []store.StreamEntry
}

func parseReadId(id string) (int64, int64, bool) {
	parts := strings.SplitN(id, "-", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}
	ms, e1 := strconv.ParseInt(parts[0], 10, 64)
	seq, e2 := strconv.ParseInt(parts[1], 10, 64)
	return ms, seq, e1 == nil && e2 == nil
}

func readStreams(keys, ids []string) []StreamData {
	var result []StreamData

	for i := range keys {
		key := keys[i]
		startMs, startSeq, ok := parseReadId(ids[i])
		if !ok {
			return nil
		}

		store.Mu.Lock()
		stream := store.Streams[key]
		store.Mu.Unlock()

		var entries []store.StreamEntry
		for _, e := range stream {
			ms, seq, _ := parseReadId(e.Id)
			if ms > startMs || (ms == startMs && seq > startSeq) {
				entries = append(entries, e)
			}
		}

		if len(entries) > 0 {
			result = append(result, StreamData{Key: key, Entries: entries})
		}
	}
	return result
}

func handleBlock(conn net.Conn, args []string) {
	timeoutMs, _ := strconv.Atoi(args[2])
	streamCount := (len(args) - 4) / 2
	keys := args[4 : 4+streamCount]
	ids := args[4+streamCount:]

	ch := make(chan struct{}, 1)

	store.Mu.Lock()
	for _, key := range keys {
		store.StreamBlockingClients[key] = append(store.StreamBlockingClients[key], ch)
	}
	store.Mu.Unlock()

	if timeoutMs == 0 {
		<-ch
		result := readStreams(keys, ids)
		if len(result) == 0 {
			resp.WriteNullArray(conn)
			return
		}

		writeXReadResponse(conn, result)
		return
	}

	select {
	case <-ch:
	case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
		resp.WriteNullArray(conn)
		return
	}

	result := readStreams(keys, ids)
	if len(result) == 0 {
		resp.WriteNullArray(conn)
		return
	}

	writeXReadResponse(conn, result)
}

func writeXReadResponse(conn net.Conn, result []StreamData) {
	fmt.Fprintf(conn, "*%d\r\n", len(result))
	for _, s := range result {
		fmt.Fprintf(conn, "*2\r\n")
		resp.WriteBulkString(conn, s.Key)
		fmt.Fprintf(conn, "*%d\r\n", len(s.Entries))
		for _, e := range s.Entries {
			fmt.Fprintf(conn, "*2\r\n")
			resp.WriteBulkString(conn, e.Id)
			fmt.Fprintf(conn, "*%d\r\n", len(e.Fields))
			for _, f := range e.Fields {
				resp.WriteBulkString(conn, f)
			}
		}
	}
}

func xreadHandler(conn net.Conn, args []string) {
	if len(args) < 4 {
		resp.WriteError(conn, "ERR: wrong number of arguments for XREAD")
		return
	}

	if strings.ToUpper(args[1]) == "BLOCK" {
		handleBlock(conn, args)
		return
	}

	streamCount := (len(args) - 2) / 2
	keys := args[2 : 2+streamCount]
	ids := args[2+streamCount:]

	result := readStreams(keys, ids)
	if len(result) == 0 {
		resp.WriteNullArray(conn)
		return
	}

	writeXReadResponse(conn, result)
}