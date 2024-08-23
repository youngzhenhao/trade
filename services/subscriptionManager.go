package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"trade/models"
)

type SubscriptionManager struct {
	subscriptions map[string]map[*models.Client]bool
	mutex         sync.RWMutex
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: make(map[string]map[*models.Client]bool),
	}
}

func (sm *SubscriptionManager) Subscribe(key string, client *models.Client) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if _, exists := sm.subscriptions[key]; !exists {
		sm.subscriptions[key] = make(map[*models.Client]bool)
	}
	sm.subscriptions[key][client] = true
	log.Printf("Client %v subscribed to key: %s. Total subscribers: %d", client, key, len(sm.subscriptions[key]))
}

func (sm *SubscriptionManager) Unsubscribe(key string, client *models.Client) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if clients, exists := sm.subscriptions[key]; exists {
		delete(clients, client)
		if len(clients) == 0 {
			delete(sm.subscriptions, key)
		}
	}
}

func (sm *SubscriptionManager) Broadcast(key string, message interface{}) {
	sm.mutex.RLock()
	clients := sm.subscriptions[key]
	sm.mutex.RUnlock()

	// 将消息转换为 JSON 字节数组
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println("Failed to marshal message:", err)
		return
	}

	// 创建一个格式化的 JSON 字符串用于打印
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, jsonMessage, "", "  ")
	if err != nil {
		log.Println("Failed to format JSON for printing:", err)
	} else {
		fmt.Println("Broadcasting message:")
		fmt.Println(prettyJSON.String())
	}

	for client := range clients {
		select {
		case client.Send <- jsonMessage:
		default:
			log.Println("Failed to send message to client, channel might be full")
		}
	}
}
