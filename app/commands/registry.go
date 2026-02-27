package commands

import "net"

type CommandFunc func(conn net.Conn, args []string)

var Registry = map[string]CommandFunc{
	"PING":   pingHandler,
	"ECHO":   echoHandler,
	"SET":    setHandler,
	"GET":    getHandler,
	"RPUSH":  rpushHandler,
	"LPUSH":  lpushHandler,
	"LRANGE": lrangeHandler,
	"LLEN":   llenHandler,
	"LPOP":   lpopHandler,
	"BLPOP":  blpopHandler,

	"TYPE":   typeHandler,
	"XADD":   xaddHandler,
}
