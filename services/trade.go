package services

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
	"trade/middleware"
	"trade/models"

	"github.com/google/uuid"
)

type TransactionService struct {
	Orders              map[string]*models.TradeOrder
	OrderStatus         map[string]bool
	Mutex               *sync.RWMutex
	Clients             map[string]*models.Client
	subscriptionManager *SubscriptionManager
	messageQueue        chan messageTask
}

type messageTask struct {
	client  *models.Client
	message []byte
}

func NewTransactionService() *TransactionService {
	ts := &TransactionService{
		Orders:              make(map[string]*models.TradeOrder),
		OrderStatus:         make(map[string]bool),
		Mutex:               &sync.RWMutex{},
		Clients:             make(map[string]*models.Client),
		subscriptionManager: NewSubscriptionManager(),
		messageQueue:        make(chan messageTask, 1000),
	}
	go ts.processMessageQueue()
	return ts
}

func (ts *TransactionService) processMessageQueue() {
	for task := range ts.messageQueue {
		ts.processMessage(task.client, task.message)
	}
}

func (ts *TransactionService) ProcessMessage(client *models.Client, msg []byte) {
	select {
	case ts.messageQueue <- messageTask{client: client, message: msg}:
		// Message added to queue
	default:
		log.Println("Message queue is full, dropping message")
	}
}

func (ts *TransactionService) SendDirectMessage(username string, message interface{}) error {
	ts.Mutex.RLock()
	client, exists := ts.Clients[username]
	ts.Mutex.RUnlock()

	if !exists {
		return fmt.Errorf("client not found: %s", username)
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	select {
	case client.Send <- jsonMessage:
		log.Printf("Sent direct message to %s: %s", username, string(jsonMessage))
	default:
		return fmt.Errorf("failed to send message to %s, channel might be full", username)
	}

	return nil
}

func (ts *TransactionService) AddClient(client *models.Client) {
	ts.Mutex.Lock()
	defer ts.Mutex.Unlock()
	ts.Clients[client.Username] = client
}

func (ts *TransactionService) RemoveClient(client *models.Client) {
	ts.Mutex.Lock()
	defer ts.Mutex.Unlock()
	delete(ts.Clients, client.Username)
}

func (ts *TransactionService) UpdateClientOrders(username string, online bool) {
	ts.Mutex.Lock()
	defer ts.Mutex.Unlock()
	for _, order := range ts.Orders {
		if order.Seller == username || order.Buyer == username {
			if online {
				order.Status = -1
			}
			ts.OrderStatus[order.OrderID] = online
		}
	}
}

func (ts *TransactionService) BroadcastMessage(message interface{}) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	ts.Mutex.RLock()
	defer ts.Mutex.RUnlock()
	for _, client := range ts.Clients {
		select {
		case client.Send <- jsonMessage:
			// Message sent successfully
		default:
			go ts.RemoveClient(client)
			close(client.Send)
			log.Printf("Failed to send message to client: %s. Closing connection.", client.Username)
		}
	}
	return nil
}

func (ts *TransactionService) GetClient(username string) (*models.Client, bool) {
	ts.Mutex.RLock()
	defer ts.Mutex.RUnlock()
	client, exists := ts.Clients[username]
	return client, exists
}

func (ts *TransactionService) processMessage(client *models.Client, msg []byte) {
	var incomingMessage map[string]interface{}
	err := json.Unmarshal(msg, &incomingMessage)
	if err != nil {
		ts.sendErrorResponse(client, "", "Invalid message format")
		return
	}

	requestID, _ := incomingMessage["request_id"].(string)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	action, ok := incomingMessage["action"].(string)
	if !ok {
		ts.sendErrorResponse(client, requestID, "Missing action")
		return
	}

	content, ok := incomingMessage["content"].(map[string]interface{})
	if !ok {
		ts.sendErrorResponse(client, requestID, "Missing content")
		return
	}

	switch action {
	case "create_order":
		ts.handleCreateOrder(client, content, requestID)
	case "query_price":
		ts.handleQueryPrice(client, content, requestID)
	case "accept_order":
		ts.handleAcceptOrder(client, content, requestID)
	case "send_psbt":
		ts.handleSendPSBT(client, content, requestID)
	case "subscribe":
		ts.handleSubscribe(client, content)
	case "unsubscribe":
		ts.handleUnsubscribe(client, content)
	case "send_direct_message":
		ts.handleSendDirectMessage(client, content, requestID)
	default:
		ts.sendErrorResponse(client, requestID, "Unknown action")
	}
}

func (ts *TransactionService) handleCreateOrder(client *models.Client, content map[string]interface{}, requestID string) {
	var order models.TradeOrder
	if err := mapToStruct(content, &order); err != nil {
		ts.sendErrorResponse(client, requestID, "Invalid order format")
		return
	}
	ts.createNewOrder(client, order, requestID)
}

func (ts *TransactionService) handleQueryPrice(client *models.Client, content map[string]interface{}, requestID string) {
	var queryParams models.QueryOrder
	if err := mapToStruct(content, &queryParams); err != nil {
		ts.sendErrorResponse(client, requestID, "Invalid query parameters format")
		return
	}
	ts.matchUnitPrice(client, queryParams, requestID)
}

func (ts *TransactionService) handleAcceptOrder(client *models.Client, content map[string]interface{}, requestID string) {
	var order models.TradeOrder
	if err := mapToStruct(content, &order); err != nil {
		ts.sendErrorResponse(client, requestID, "Invalid order format")
		return
	}
	ts.updateOrder(client, order, requestID)
}

func (ts *TransactionService) handleSendPSBT(client *models.Client, content map[string]interface{}, requestID string) {
	var order models.TradeOrder
	if err := mapToStruct(content, &order); err != nil {
		ts.sendErrorResponse(client, requestID, "Invalid order format")
		return
	}
	ts.confirmPsbtSignInfo(client, order, requestID)
}

func (ts *TransactionService) handleSubscribe(client *models.Client, content map[string]interface{}) {
	channel, ok := content["channel"].(string)
	if !ok {
		ts.sendErrorResponse(client, "", "Invalid subscription channel")
		return
	}

	interval, ok := content["interval"].(float64)
	if !ok {
		ts.sendErrorResponse(client, "", "Invalid interval")
		return
	}
	tradingPair, ok := content["tradingPair"].(string)
	if !ok {
		ts.sendErrorResponse(client, "", "Invalid symbol")
		return
	}
	Online, ok := content["Online"]
	if !ok {
		ts.sendErrorResponse(client, "", "Invalid online")
		return
	}
	var onlineType string
	if Online == true {
		onlineType = "0"
	} else {
		onlineType = "1"
	}
	orderType, ok := content["OrderType"].(string)
	if !ok {
		ts.sendErrorResponse(client, "", "Invalid OrderType")
		return
	}
	status, ok := content["Status"].(string)
	if !ok {
		ts.sendErrorResponse(client, "", "Invalid Status")
		return
	}

	switch channel {
	case "subscriptionChannelTicker":
		subscriptionKey := fmt.Sprintf("%s:%s:%s:%s:%s", channel, tradingPair, onlineType, orderType, status)
		ts.subscriptionManager.Subscribe(subscriptionKey, client)

		if len(ts.subscriptionManager.subscriptions[subscriptionKey]) == 1 {
			go ts.startDataCollection(subscriptionKey, time.Duration(interval)*time.Millisecond)
		}

		ts.sendResponse(client, map[string]interface{}{
			"action":  "subscribe",
			"status":  200,
			"message": "Subscribed to " + channel,
			"code":    200,
		})
	default:
		ts.sendErrorResponse(client, "", "Unknown subscription channel")
	}
}

func (ts *TransactionService) handleUnsubscribe(client *models.Client, content map[string]interface{}) {
	channel, ok := content["channel"].(string)
	if !ok {
		ts.sendErrorResponse(client, "", "Invalid unsubscription channel")
		return
	}

	tradingPair, ok := content["tradingPair"].(string)
	if !ok {
		ts.sendErrorResponse(client, "", "Invalid symbol")
		return
	}
	Online, ok := content["Online"]
	if !ok {
		ts.sendErrorResponse(client, "", "Invalid online")
		return
	}
	var onlineType string
	if Online == true {
		onlineType = "0"
	} else {
		onlineType = "1"
	}
	orderType, ok := content["OrderType"].(string)
	if !ok {
		ts.sendErrorResponse(client, "", "Invalid OrderType")
		return
	}
	status, ok := content["Status"].(string)
	if !ok {
		ts.sendErrorResponse(client, "", "Invalid Status")
		return
	}
	subscriptionKey := fmt.Sprintf("%s:%s:%s:%s:%s", channel, tradingPair, onlineType, orderType, status)
	ts.subscriptionManager.Unsubscribe(subscriptionKey, client)
	ts.sendResponse(client, map[string]interface{}{
		"action":  "unsubscribe",
		"status":  200,
		"message": "Unsubscribed from " + channel,
		"code":    200,
	})
}

func (ts *TransactionService) handleSendDirectMessage(client *models.Client, content map[string]interface{}, requestID string) {
	recipient, ok := content["recipient"].(string)
	if !ok {
		ts.sendErrorResponse(client, requestID, "Invalid recipient")
		return
	}

	message, ok := content["message"]
	if !ok {
		ts.sendErrorResponse(client, requestID, "Invalid message")
		return
	}

	directMessage := map[string]interface{}{
		"action":     "direct_message",
		"from":       client.Username,
		"message":    message,
		"request_id": requestID,
	}

	err := ts.SendDirectMessage(recipient, directMessage)

	if err != nil {
		ts.sendErrorResponse(client, requestID, fmt.Sprintf("Failed to send direct message: %v", err))
		return
	}

	ts.sendResponse(client, map[string]interface{}{
		"request_id": requestID,
		"status":     200,
		"message":    "Direct message sent successfully",
		"code":       200,
	})
}

func (ts *TransactionService) startDataCollection(subscriptionKey string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ts.Mutex.RLock()
			if len(ts.subscriptionManager.subscriptions[subscriptionKey]) == 0 {
				ts.Mutex.RUnlock()
				return // Stop if there are no more subscribers
			}
			data := ts.collectData(subscriptionKey)
			ts.Mutex.RUnlock()

			if data != nil {
				ts.subscriptionManager.Broadcast(subscriptionKey, data)
			}
		}
	}
}

func (ts *TransactionService) collectData(subscriptionKey string) interface{} {
	parts := strings.Split(subscriptionKey, ":")
	if len(parts) != 5 {
		return nil
	}
	channel, tradingPair, onlineType, orderType, statusStr := parts[0], parts[1], parts[2], parts[3], parts[4]
	var online bool
	if onlineType == "0" {
		online = true
	} else {
		online = false
	}

	status, err := strconv.ParseInt(statusStr, 10, 16)
	if err != nil {
		return nil
	}

	orders, err := ts.findOrdersByBaseCondition(online, tradingPair, orderType, int16(status))
	if err != nil {
		return nil
	}

	return map[string]interface{}{
		"action":      channel,
		"content":     orders,
		"tradingPair": tradingPair,
	}
}

func (ts *TransactionService) createNewOrder(client *models.Client, order models.TradeOrder, requestID string) {
	ts.Mutex.Lock()
	defer ts.Mutex.Unlock()

	orderID := ts.generateOrderID()
	order.OrderID = orderID
	order.Status = 0
	if order.OrderType == "buy" {
		order.Buyer = client.Username
	} else {
		order.Seller = client.Username
	}
	ts.Orders[orderID] = &order
	err := middleware.DB.Create(&order).Error
	if err != nil {
		ts.sendErrorResponse(client, requestID, "Failed to save order to database")
		return
	}

	response := map[string]interface{}{
		"request_id": requestID,
		"status":     200,
		"order_id":   orderID,
		"message":    "Order created successfully",
		"code":       200,
	}
	ts.sendResponse(client, response)
}

func (ts *TransactionService) matchUnitPrice(client *models.Client, queryParams models.QueryOrder, requestID string) {
	priceRange, err := ts.findOrdersByUnitPriceRange(queryParams.MinUnitPrice, queryParams.MaxUnitPrice, queryParams.Online, queryParams.TradingPair, queryParams.OrderType)
	if err != nil || len(priceRange) == 0 {
		ts.sendErrorResponse(client, requestID, "Not match unit price")
		return
	}

	priceRangeJSON, err := json.Marshal(priceRange)
	if err != nil {
		ts.sendErrorResponse(client, requestID, "Error processing matched orders")
		return
	}

	response := map[string]interface{}{
		"request_id":  requestID,
		"status":      200,
		"price_range": json.RawMessage(priceRangeJSON),
		"message":     "Match unit price",
		"code":        200,
	}
	ts.sendResponse(client, response)
}

func (ts *TransactionService) updateOrder(client *models.Client, order models.TradeOrder, requestID string) {
	ts.Mutex.Lock()
	defer ts.Mutex.Unlock()

	var existingOrder models.TradeOrder
	err := middleware.DB.Where("order_id = ?", order.OrderID).First(&existingOrder).Error
	if err != nil {
		ts.sendErrorResponse(client, requestID, "Order not found")
		return
	}
	existingOrder.Status = 1
	ts.Orders[order.OrderID] = &existingOrder
	if existingOrder.OrderType == "buy" {
		existingOrder.Seller = client.Username
	} else {
		existingOrder.Buyer = client.Username
	}

	err = middleware.DB.Save(&existingOrder).Error
	if err != nil {
		ts.sendErrorResponse(client, requestID, "Failed to update order in database")
		return
	}

	response := map[string]interface{}{
		"request_id": requestID,
		"status":     200,
		"message":    "Order updated successfully",
		"code":       200,
	}
	ts.sendResponse(client, response)
}

func (ts *TransactionService) confirmPsbtSignInfo(client *models.Client, order models.TradeOrder, requestID string) {
	ts.Mutex.Lock()
	defer ts.Mutex.Unlock()

	var existingOrder models.TradeOrder
	err := middleware.DB.Where("order_id = ?", order.OrderID).First(&existingOrder).Error
	if err != nil {
		ts.sendErrorResponse(client, requestID, "Order not found")
		return
	}
	if order.OrderType == "buy" {
		if existingOrder.Buyer != client.Username {
			ts.sendErrorResponse(client, requestID, "Unauthorized to update order")
			return
		}
		existingOrder.PSBTBuyer = order.PSBTBuyer
	} else {
		if existingOrder.Seller != client.Username {
			ts.sendErrorResponse(client, requestID, "Unauthorized to update order")
			return
		}
		existingOrder.PSBTSeller = order.PSBTSeller
	}
	err = middleware.DB.Save(&existingOrder).Error
	if err != nil {
		ts.sendErrorResponse(client, requestID, "Failed to update order in database")
		return
	}
	ts.Orders[order.OrderID] = &existingOrder
	response := map[string]interface{}{
		"request_id": requestID,
		"status":     200,
		"message":    "Order updated successfully",
		"code":       200,
	}
	ts.sendResponse(client, response)
}

func (ts *TransactionService) findOrdersByBaseCondition(online bool, tradingPair, orderType string, status int16) ([]models.TradeOrder, error) {
	var orders []models.TradeOrder
	err := middleware.DB.Where("online = ? AND trading_pair = ? AND order_type = ? and status=?", online, tradingPair, orderType, status).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (ts *TransactionService) findOrdersByUnitPriceRange(minPrice, maxPrice float64, online bool, tradingPair, orderType string) ([]models.TradeOrder, error) {
	var orders []models.TradeOrder
	err := middleware.DB.Where("unit_price >= ? AND unit_price <= ? AND online = ? AND trading_pair = ? AND order_type = ?", minPrice, maxPrice, online, tradingPair, orderType).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (ts *TransactionService) sendErrorResponse(client *models.Client, requestID, message string) {
	response := map[string]interface{}{
		"request_id": requestID,
		"status":     -999,
		"message":    message,
		"code":       -999,
	}
	ts.sendResponse(client, response)
}

func (ts *TransactionService) sendResponse(client *models.Client, response map[string]interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling response:", err)
		return
	}
	select {
	case client.Send <- jsonResponse:
		log.Println("Sent message to client:", string(jsonResponse))
	default:
		log.Println("Client send channel is blocked, discarding message")
	}
}

func (ts *TransactionService) generateOrderID() string {
	return strconv.Itoa(rand.Int())
}

func mapToStruct(input map[string]interface{}, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}

func (ts *TransactionService) CloseAllConnections() {
	ts.Mutex.Lock()
	defer ts.Mutex.Unlock()
	for _, client := range ts.Clients {
		client.Conn.Close()
		close(client.Send)
	}
	ts.Clients = make(map[string]*models.Client)
}
