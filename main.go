package main

import (
	"fmt"
	"net" // net package to open standard TCP sockets
	"os"
	"strings"
)

type Handler func(conn net.Conn, method string)

// maps string paths to their respective handler functions
type Router struct {
	routes map[string]Handler
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]Handler)}
}

func (r *Router) Handle(path string, handler Handler) {
	r.routes[path] = handler
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Failed to bind: %s\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server is listening on http://localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		buffer := make([]byte, 1024)
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			conn.Close()
			continue
		}

		request := string(buffer[:bytesRead])

		lines := strings.Split(request, "\n")

		var method, path string

		if len(lines) > 0 {
			requestLine := strings.TrimSpace(lines[0])

			parts := strings.Split(requestLine, " ")

			if len(parts) >= 3 {
				method = parts[0]
				path = parts[1]

				fmt.Printf("Received %s request for path: %s\n", method, path)
			}
		}

		body := "<html><body><h1>Wello Horld!</h1><p>No usage of <code>net/http</code></p></body></html>"

		response := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"Content-Length: %d\r\n"+
			"Connection: close\r\n"+
			"\r\n"+
			"%s", len(body), body)

		conn.Write([]byte(response))

		conn.Close()
	}
}
