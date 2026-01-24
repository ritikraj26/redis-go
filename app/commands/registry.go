package commands

import "net"

type CommandFunc func(conn net.Conn, args []string)

var Registry = map[string]CommandFunc{
	"PING":  Ping,
	"ECHO":  Echo,
	"SET":   Set,
	"GET":   Get,
	"RPUSH": RPush,
}
