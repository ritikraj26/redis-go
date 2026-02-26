package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

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
			resp.WriteError(conn, "Error reading input")
			return
		}

		if len(args) == 0 {
			continue
		}

		for i := 0; i < len(args); i++ {
			fmt.Println(args[i])
		}

		cmd := strings.ToUpper(args[0])

		if handler, ok := commands.Registry[cmd]; ok {
			handler(conn, args)
		} else {
			resp.WriteError(conn, "Unknown command: "+cmd)
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
