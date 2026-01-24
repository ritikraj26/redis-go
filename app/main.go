package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

// var data = make(map[string]string)
type data struct {
	value  string
	expiry *time.Time
}

var list = make(map[string][]string)

var store = make(map[string]data)

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
			if len(args) < 3 {
				writeError(conn, "SET requires arguments (key and value)")
			} else {
				key := args[1]
				value := args[2]

				var expiry *time.Time

				if len(args) >= 5 {
					ttl, err := strconv.Atoi(args[4])
					if err != nil {
						writeError(conn, "ERR invalid expire time")
						return
					}

					var d time.Duration
					switch strings.ToUpper(args[3]) {
					case "EX":
						d = time.Duration(ttl) * time.Second
					case "PX":
						d = time.Duration(ttl) * time.Millisecond
					default:
						writeError(conn, "ERR invalid expire option")
						return
					}

					t := time.Now().Add(d)
					expiry = &t
				}

				store[key] = data{
					value:  value,
					expiry: expiry, // nil if not provided
				}
				writeSimpleString(conn, "OK")
			}
		case "GET":
			if len(args) < 2 {
				writeError(conn, "Get requires an argument (key)")
			} else {
				val, ok := store[args[1]]
				if ok {
					if val.expiry != nil {
						if time.Now().After(*val.expiry) {
							delete(store, args[1])
							writeNullString(conn)
						} else {
							writeBulkString(conn, val.value)
						}
					} else {
						writeBulkString(conn, val.value)
					}

				} else {
					writeError(conn, "Key not present")
				}
			}
		case "RPUSH":
			if len(args) < 3 {
				writeError(conn, "Too few arguments for RPUSH")
			} else {
				val, ok := list[args[1]]
				if ok {
					for i := 2; i < len(args); i++ {
						val = append(val, args[i])
					}
					list[args[1]] = val
					// writeSimpleString(conn, strconv.Itoa(len(val)))
					writeInteger(conn, uint32(len(val)))
				} else {
					val = []string{}
					for i := 2; i < len(args); i++ {
						val = append(val, args[i])
					}
					list[args[1]] = val
					writeInteger(conn, uint32(len(val)))
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

func writeInteger(conn net.Conn, value uint32) {
	conn.Write([]byte(fmt.Sprintf(":%d\r\n", value)))
}

func writeSimpleString(conn net.Conn, message string) {
	conn.Write([]byte("+" + message + "\r\n"))
}

func writeBulkString(conn net.Conn, message string) {
	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)))
}

func writeNullString(conn net.Conn) {
	conn.Write([]byte("$-1\r\n"))
}

func writeError(conn net.Conn, message string) {
	conn.Write([]byte("-" + message + "\r\n"))
}
