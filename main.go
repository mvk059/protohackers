package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000" // Default port if not specified
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	log.Printf("Server listening on 0.0.0.0:%s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	_, err := io.Copy(conn, conn)
	if err != nil {
		log.Printf("Failed to copy: %v", err)
		return
	}
}
