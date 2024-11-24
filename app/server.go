package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed after handling

	req := make([]byte, 1024)
	if _, err := conn.Read(req); err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	// Extract the request route
	requestLine := strings.SplitN(string(req), "\r\n", 2)[0]
	parts := strings.Fields(requestLine) // Splits "GET / HTTP/1.1" into ["GET", "/", "HTTP/1.1"]
	if len(parts) < 2 {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	route := parts[1] // Extract the route

	// Handle routes
	switch {
	case route == "/":
		sendResponse(conn, "200 OK", "Connection Established!")
	case strings.HasPrefix(route, "/echo/"):
		param := strings.TrimPrefix(route, "/echo/")
		sendResponse(conn, "200 OK", param)
	default:
		sendResponse(conn, "404 Not Found", "Route not found")
	}
}

func sendResponse(conn net.Conn, status, body string) {
	response := fmt.Sprintf(
		"HTTP/1.1 %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		status, len(body), body,
	)
	conn.Write([]byte(response))
}

func main() {
	port := "4221"
	fmt.Printf("Starting server on port %s...\n", port)

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Printf("Error starting server on port %s: %s\n", port, err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn) // Handle connections concurrently
	}
}
