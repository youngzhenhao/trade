package models

import "github.com/gorilla/websocket"

type Client struct {
	Conn          *websocket.Conn
	Send          chan []byte
	Username      string
	Done          chan struct{}
	Subscriptions map[string]Subscription
	Ack           chan struct{}
}
