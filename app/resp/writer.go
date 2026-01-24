package resp

import (
	"fmt"
	"net"
)

func WriteInteger(conn net.Conn, value uint32) {
	fmt.Fprintf(conn, ":%d\r\n", value)
}

func WriteSimpleString(conn net.Conn, message string) {
	fmt.Fprintf(conn, "+%s\r\n", message)
}

func WriteBulkString(conn net.Conn, message string) {
	fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(message), message)
}

func WriteNullString(conn net.Conn) {
	fmt.Fprint(conn, "$-1\r\n")
}

func WriteError(conn net.Conn, message string) {
	fmt.Fprintf(conn, "-%s\r\n", message)
}
