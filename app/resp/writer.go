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

func WriteNullBulkString(conn net.Conn) {
	fmt.Fprint(conn, "$-1\r\n")
}

func WriteError(conn net.Conn, message string) {
	fmt.Fprintf(conn, "-%s\r\n", message)
}

func WriteBulkStringArray(conn net.Conn, elements []string) {
	fmt.Fprintf(conn, "*%d\r\n", len(elements))
	for _, element := range elements {
		WriteBulkString(conn, element)
	}
}

func WriteNullArray(conn net.Conn) {
	fmt.Fprint(conn, "*-1\r\n")
}