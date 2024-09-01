package utils

import (
	"math"
	"net"
)

func SendMalformedResponse(conn net.Conn) {
	malformedResp := []byte("{\"method\":\"isPrime\",\"prime\":}\n")
	conn.Write(malformedResp)
}

func IsNumberPrime(num float64) bool {
	if num <= 1 || math.Floor(num) != num {
		return false
	}
	if num == 2 {
		return true
	}
	if math.Mod(num, 2) == 0 {
		return false
	}
	for i := 3.0; i <= math.Sqrt(num); i += 2 {
		if math.Mod(num, i) == 0 {
			return false
		}
	}
	return true
}
