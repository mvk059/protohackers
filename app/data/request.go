package data

import "encoding/json"

type Request struct {
	Method string      `json:"method"`
	Number json.Number `json:"number"`
}
