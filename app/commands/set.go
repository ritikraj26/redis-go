package commands

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func setHandler(conn net.Conn, args []string) {
	if len(args) < 3 {
		resp.WriteError(conn, "SET requires arguments (key and value)")
		return
	}

	key := args[1]
	value := args[2]

	var expiry *time.Time

	if len(args) >= 5 {
		ttl, err := strconv.Atoi(args[4])
		if err != nil {
			resp.WriteError(conn, "ERR invalid expire time")
			return
		}

		var d time.Duration
		switch strings.ToUpper(args[3]) {
		case "EX":
			d = time.Duration(ttl) * time.Second
		case "PX":
			d = time.Duration(ttl) * time.Millisecond
		default:
			resp.WriteError(conn, "ERR invalid expire option")
			return
		}

		t := time.Now().Add(d)
		expiry = &t
	}

	store.Store[key] = store.Data{
		Value:  value,
		Expiry: expiry,
	}

	resp.WriteSimpleString(conn, "OK")
}
