package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

var data = make(map[string]string)

func main() {
	// Start listening on TCP port 6379 on all interfaces.
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		// Handle the connection in a new goroutine.
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close() // execute at the end of the function

	// Create a buffered reader to read lines from the connection.
	reader := bufio.NewReader(conn)
	fmt.Println("New client connected:", conn.RemoteAddr().String())
	for {
		// Read a line from the connection.
		args, err := readRESP(reader)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			writeError(conn, "Error reading input")
			return
		}

		if len(args) == 0 {
			continue
		}

		for i := 0; i < len(args); i++ {
			fmt.Println(args[i])
		}

		cmd := strings.ToUpper(args[0])

		switch cmd {
		case "PING":
			writeSimpleString(conn, "PONG")
		case "ECHO":
			if len(args) < 2 {
				writeError(conn, "ECHO requires an argument")
			} else {
				writeBulkString(conn, args[1])
			}
		case "SET":
			if len(args) < 2 {
				writeError(conn, "SET requires arguments (key and value)")
			} else {
				data[args[1]] = args[2]
				writeSimpleString(conn, "OK")
			}
		case "GET":
			if len(args) < 2 {
				writeError(conn, "Get requires an argument (key)")
			} else {
				val, ok := data[args[1]]
				if ok {
					writeBulkString(conn, val)
				} else {
					writeError(conn, "Key not present")
				}
			}
		default:
			writeError(conn, "Unknown command: "+cmd)
		}
	}
}

func readRESP(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] != '*' {
		return nil, nil
	}

	numArgs, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, err
	}

	args := make([]string, 0, numArgs)
	for i := 0; i < numArgs; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] != '$' {
			return nil, nil
		}

		argLen, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}

		arg := make([]byte, argLen+2) // +2 for \r\n
		_, err = reader.Read(arg)
		fmt.Println("Read arg:", string(arg))
		if err != nil {
			return nil, err
		}

		args = append(args, string(arg[:argLen]))
	}
	fmt.Println("Parsed args:", args)
	return args, nil
}

func writeSimpleString(conn net.Conn, message string) {
	conn.Write([]byte("+" + message + "\r\n"))
}

func writeBulkString(conn net.Conn, message string) {
	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)))
}

func writeError(conn net.Conn, message string) {
	conn.Write([]byte("-" + message + "\r\n"))
}
