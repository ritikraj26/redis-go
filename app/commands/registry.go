package commands

import "net"

type CommandFunc func(conn net.Conn, args []string)

var Registry = map[string]CommandFunc{
	"PING":   pingHandler,
	"ECHO":   echoHandler,
	"SET":    setHandler,
	"GET":    getHandler,
	"RPUSH":  rpushHandler,
	"LRANGE": lrangeHandler,
}
