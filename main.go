package main

import (
	"fmt"
	"net" // net package to open standard TCP sockets
	"os"
)

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

		fmt.Println("--- REQUEST STARTS ---")
		fmt.Print(string(buffer[:bytesRead]))
		fmt.Println("--- REQUEST ENDS ---")

		conn.Close()
	}
}
