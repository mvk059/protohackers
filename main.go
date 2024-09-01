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
		log.Printf("Received request: %s", line)

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Printf("Malformed JSON: %v", err)
			sendMalformedResponse(conn)
			return
		}

		if req.Method != "isPrime" {
			log.Printf("Invalid method: %s", req.Method)
			sendMalformedResponse(conn)
			return
		}

		if req.Number == "" {
			log.Printf("Missing number field")
			sendMalformedResponse(conn)
			return
		}

		num, err := req.Number.Float64()
		if err != nil {
			log.Printf("Invalid number: %v", err)
			sendMalformedResponse(conn)
			return
		}

		isPrime := isPrime(int(num))
		resp := Response{Method: "isPrime", Prime: isPrime}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			return
		}

		jsonResp = append(jsonResp, '\n')
		log.Printf("Sending response: %s", string(jsonResp))
		_, err = conn.Write(jsonResp)
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
	log.Printf("Sending malformed response: %s", string(malformedResp))
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
