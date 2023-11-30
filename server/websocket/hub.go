package websocket

import (
	"time"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/infras/logger"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	loginAt time.Time
}

type MessageType int

const (
	MessageTypeInvalid MessageType = iota
	MessageTypeBroadcast
	MessageTypeGroup
	MessageTypePrivate
)

type Message struct {
	Type     MessageType `json:"type"`
	Metadata string      `json:"metadata"`

	From      string    `json:"from"` // user subject
	To        string    `json:"to"`   // user subject
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func NewHub(env environment.Env, logger logger.Logger) (*Hub, error) {
	return &Hub{
		logger:  logger,
		clients: make(map[string]*Client),
	}, nil
}

type Hub struct {
	logger  logger.Logger
	clients map[string]*Client
}

func (h *Hub) RegisterClient(userSubject string, conn *websocket.Conn) error {
	client, ok := h.clients[userSubject]
	if ok {
		client.conn.Close()
		delete(h.clients, userSubject)
	}

	client = &Client{
		conn:    conn,
		loginAt: time.Now(),
	}

	h.clients[userSubject] = client
	return nil
}

func (h *Hub) UnregisterClient(key string) error {
	client, ok := h.clients[key]
	if !ok {
		return nil
	}

	client.conn.Close()
	delete(h.clients, key)
	return nil
}

func (h *Hub) SendMessage(m Message) error {
	client, ok := h.clients[m.To]
	if !ok {
		// TODO: add to storage
		return nil
	}

	if err := client.conn.WriteJSON(m); err != nil {
		// TODO: into message queue
		h.logger.Errorf("消息[%v]发送失败：%w", m, err)
	}

	return nil
}

func (h *Hub) Close() error {
	for _, c := range h.clients {
		c.conn.Close()
	}

	return nil
}
