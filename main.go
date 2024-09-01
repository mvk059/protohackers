package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
)

type Request struct {
	Method string      `json:"method"`
	Number json.Number `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

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
		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil || req.Method != "isPrime" || req.Number == "" {
			sendMalformedResponse(conn)
			log.Printf("Malformed request: %v", line)
			return
		}

		num, err := req.Number.Float64()
		if err != nil {
			sendMalformedResponse(conn)
			log.Printf("Malformed request number: %v", line)
			return
		}

		isPrime := isPrime(int(num))
		response := Response{
			Method: "isPrime",
			Prime:  isPrime,
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			return
		}

		log.Printf("Response: %v for request: %v", response, line)

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

func sendMalformedResponse(conn net.Conn) {
	malformedResp := []byte("{\"method\":\"isPrime\",\"prime\":}\n")
	conn.Write(malformedResp)
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	limit := int(math.Sqrt(float64(n)))
	for i := 5; i <= limit; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}
