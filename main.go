package main

import (
	"fmt"
	"net" // net package to open standard TCP sockets
	"os"
	"strings"
	"time"
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

	router := NewRouter()
	router.Handle("/", handleHome)
	router.Handle("/weirdshit", handleSome)

	fmt.Println("Server is listening on http://localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// 'go' keyword spawns a lightweight thread (Goroutine)
		go handleClient(conn, router)
	}
}

func handleHome(conn net.Conn, method string) {
	body := "<html><body><h1>Welcome to the HomePage!</h1><p>Served concurrently via Goroutines</p></body></html>"
	sendHTML(conn, "200 OK", body)
}

func handleSome(conn net.Conn, method string) {
	time.Sleep(2 * time.Second)

	body := "<html><body><h1>Some weirdShit endpoint</h1><p>This request took 2 seconds!</p></body></html>"
	sendHTML(conn, "200 OK", body)
}

func sendHTML(conn net.Conn, status string, body string) {
	response := fmt.Sprintf("HTTP/1.1 %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"Content-Length: %d\r\n"+
		"Connection: close\r\n"+
		"\r\n"+
		"%s", status, len(body), body)

	conn.Write([]byte(response))
}

func handleClient(conn net.Conn, router *Router) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		return
	}

	request := string(buffer[:bytesRead])
	lines := strings.Split(request, "\n")
	if len(lines) == 0 {
		return
	}

	requestLine := strings.TrimSpace(lines[0])
	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		return
	}

	method := parts[0]
	path := parts[1]

	fmt.Printf("Handling %s %s\n", method, path)

	handler, exists := router.routes[path]
	if exists {
		handler(conn, method)
	} else {
		handleNotFound(conn)
	}
}

func handleNotFound(conn net.Conn) {
	body := "<html><body><h1>404 Not Found</h1><p>The requested path does not exist.</p></body></html>"
	sendHTML(conn, "404 Not Found", body)
}
