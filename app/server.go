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
	requestParts := strings.SplitN(string(req), "\r\n", 2) // Separate request line and headers
	requestLine := requestParts[0]                         // First line: e.g., "GET / HTTP/1.1"
	parts := strings.Fields(requestLine)                   // Split into ["GET", "/route", "HTTP/1.1"]
	if len(parts) < 2 {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	route := parts[1] // Extract the route

	// Handle routes
	switch {
	case route == "/":
		sendResponse(conn, "200 OK", "text/plain", "Connection Established!")
	case strings.HasPrefix(route, "/echo/"):
		param := strings.TrimPrefix(route, "/echo/")
		sendResponse(conn, "200 OK", "text/plain", param)
	case route == "/user-agent":
		userAgent := extractHeaderValue(string(req), "User-Agent")
		if userAgent == "" {
			sendResponse(conn, "400 Bad Request", "text/plain", "User-Agent header missing")
		} else {
			sendResponse(conn, "200 OK", "text/plain", userAgent)
		}
	case strings.HasPrefix(route, "/files/"):
		directory := os.Args[2]
		fileName := strings.TrimPrefix(route, "/files/")
		fmt.Print(fileName)
		data, err := os.ReadFile(directory + fileName)
		if err != nil {
			sendResponse(conn, "404 Not Found", "text/plain", "File not found")
		} else {
			sendResponse(conn, "200 OK", "application/octet-stream", string(data))
		}
	default:
		sendResponse(conn, "404 Not Found", "text/plain", "Route not found")
	}
}

func extractHeaderValue(request, headerName string) string {
	// Split headers from the body
	parts := strings.Split(request, "\r\n\r\n")
	if len(parts) < 1 {
		return ""
	}

	// Extract headers
	headers := strings.Split(parts[0], "\r\n")
	for _, header := range headers {
		if strings.HasPrefix(header, headerName+":") {
			return strings.TrimSpace(strings.TrimPrefix(header, headerName+":"))
		}
	}
	return ""
}

func sendResponse(conn net.Conn, status, contentType, body string) {
	response := fmt.Sprintf(
		"HTTP/1.1 %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
		status, contentType, len(body), body,
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
