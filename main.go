package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"protohackers/app/data"
	"protohackers/app/utils"
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
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		var req data.Request
		if err := json.Unmarshal([]byte(line), &req); err != nil || req.Method != "isPrime" || req.Number == "" {
			utils.SendMalformedResponse(conn)
			return
		}

		num, err := req.Number.Float64()
		if err != nil {
			utils.SendMalformedResponse(conn)
			return
		}

		isPrime := utils.IsNumberPrime(num)
		response := data.Response{
			Method: "isPrime",
			Prime:  isPrime,
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			return
		}

		jsonResponse = append(jsonResponse, '\n')
		_, err = conn.Write(jsonResponse)
		if err != nil {
			log.Printf("Failed to write response: %v", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from connection: %v", err)
	}
}
