package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func connect(conn net.Conn) {
	req := make([]byte, 1024)
	_, err := conn.Read(req)
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	if !strings.HasPrefix(string(req), "GET / HTTP/1.1") {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}

	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nConnction Established!\r\n"))
	conn.Close()
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	connect(conn)

}
