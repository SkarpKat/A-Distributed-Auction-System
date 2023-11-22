package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// Define the address and port for the server node
	address := "localhost"
	port := "8080"

	// Create a TCP listener on the specified address and port
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}

	fmt.Println("Server node started on", address+":"+port)

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// TODO: Implement your logic for handling the connection here

	// Close the connection when done
	defer conn.Close()
}
