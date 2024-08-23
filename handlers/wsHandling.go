package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"
	"trade/middleware"
	"trade/models"
	"trade/services"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 100 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024 // 1 MB
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true // 注意：在生产环境中应该更严格地检查origin
	},
}

type messageTask struct {
	client  *models.Client
	message []byte
}

var (
	transactionService *services.TransactionService
	once               sync.Once
	shutdownSignal     = make(chan struct{})
	messageQueue       = make(chan messageTask, 1000)
)

func GetTransactionService() *services.TransactionService {
	once.Do(func() {
		transactionService = services.NewTransactionService()
		go processMessageQueue()
	})
	return transactionService
}

func processMessageQueue() {
	for {
		select {
		case task := <-messageQueue:
			transactionService.ProcessMessage(task.client, task.message)
		case <-shutdownSignal:
			return
		}
	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := middleware.ValidateToken(token[7:])
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}

	client := &models.Client{
		Conn:          conn,
		Send:          make(chan []byte, 256),
		Username:      claims.Username,
		Done:          make(chan struct{}),
		Subscriptions: make(map[string]models.Subscription),
	}

	ts := GetTransactionService()
	ts.AddClient(client)

	go handleConnection(client)
	go handleMessages(client)
}

func handleConnection(client *models.Client) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in handleConnection: %v", r)
		}
		cleanupClient(client)
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-client.Done:
			return
		case <-shutdownSignal:
			return
		default:
			_, msg, err := client.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Error reading message: %v", err)
				}
				return
			}
			select {
			case messageQueue <- messageTask{client: client, message: msg}:
			default:
				log.Println("Message queue is full, dropping message")
			}
		}
	}
}

func handleMessages(client *models.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case <-client.Done:
			return
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func cleanupClient(client *models.Client) {
	transactionService.RemoveClient(client)
	transactionService.UpdateClientOrders(client.Username, false)
	client.Conn.Close()
	close(client.Done)
}

func BroadcastMessage(msg []byte) {
	transactionService.BroadcastMessage(msg)
}

func ShutdownWebSockets() {
	close(shutdownSignal)
	transactionService.CloseAllConnections()
	time.Sleep(2 * time.Second)
}
