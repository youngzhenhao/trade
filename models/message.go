package models

import "github.com/goccy/go-json"

type Message struct {
	Username string          `json:"username"`
	Action   string          `json:"action"`
	Content  json.RawMessage `json:"content"`
}
