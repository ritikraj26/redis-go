package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func pingHandler(conn net.Conn, args []string) {
	resp.WriteSimpleString(conn, "PONG")
}
