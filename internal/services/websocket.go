package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Piyu-Pika/godzilla-go/internal/database"
	// "github.com/Piyu-Pika/godzilla-go/internal/handlers"
	"github.com/Piyu-Pika/godzilla-go/internal/models"
	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	Clients    map[uint]*Client
	Broadcast  chan []byte
	SendToUser chan struct {
		UserID  uint
		Message []byte
	}
	register   chan *Client
	unregister chan *Client
	Mutex      sync.RWMutex
}

type Client struct {
	socket *websocket.Conn
	send   chan []byte
	UserID uint
}

type ChatMessage struct {
	RecipientID uint   `json:"recipient_id"`
	Content     string `json:"content"`
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Clients:   make(map[uint]*Client),
		Broadcast: make(chan []byte),
		SendToUser: make(chan struct {
			UserID  uint
			Message []byte
		}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (manager *WebSocketManager) Run() {
	for {
		select {
		case client := <-manager.register:
			manager.Mutex.Lock()
			manager.Clients[client.UserID] = client
			manager.Mutex.Unlock()
			log.Println("Registered client for user", client.UserID)
		case client := <-manager.unregister:
			manager.Mutex.Lock()
			if _, ok := manager.Clients[client.UserID]; ok {
				delete(manager.Clients, client.UserID)
				close(client.send)
			}
			manager.Mutex.Unlock()
		case message := <-manager.Broadcast:
			manager.Mutex.RLock()
			for _, client := range manager.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(manager.Clients, client.UserID)
				}
			}
			manager.Mutex.RUnlock()
		case sendData := <-manager.SendToUser:
			manager.Mutex.RLock()
			if client, ok := manager.Clients[sendData.UserID]; ok {
				select {
				case client.send <- sendData.Message:
				default:
					close(client.send)
					delete(manager.Clients, client.UserID)
				}
			}
			manager.Mutex.RUnlock()
		}
	}
}

func (manager *WebSocketManager) HandleWebSocketConnections(w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	client := &Client{
		socket: c,
		send:   make(chan []byte, 256),
		// UserID: userID,
	}
	manager.register <- client
	go manager.writePump(client)
	go manager.readPump(client)
}

func (manager *WebSocketManager) readPump(client *Client) {
	defer func() {
		manager.unregister <- client
		client.socket.Close()
	}()
	for {
		_, message, err := client.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var chatMsg ChatMessage
		if err := json.Unmarshal(message, &chatMsg); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}
		// Store message in MongoDB
		msgToStore := models.ChatMessage{
			SenderID:    client.UserID,
			RecipientID: chatMsg.RecipientID,
			Content:     chatMsg.Content,
			Timestamp:   time.Now(),
		}
		_, err = database.ChatCollection.InsertOne(context.TODO(), msgToStore)
		if err != nil {
			log.Println("Failed to store message:", err)
		}
		// Send to recipient
		manager.SendToUser <- struct {
			UserID  uint
			Message []byte
		}{chatMsg.RecipientID, message}
	}
}

func (manager *WebSocketManager) writePump(client *Client) {
	ticker := time.NewTicker(time.Second)
	defer func() {
		ticker.Stop()
		manager.unregister <- client
		client.socket.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				return
			}
			client.socket.SetWriteDeadline(time.Now().Add(time.Second * 5))
			if err := client.socket.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := client.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
