package hub

import (
	"context"
	"fangaoxs.com/go-chat/internal/domain/record"
	"time"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"github.com/gorilla/websocket"
)

type MessageType = entity.RecordType

const (
	MessageTypeInvalid   MessageType = entity.RecordTypeInvalid
	MessageTypeBroadcast MessageType = entity.RecordTypeBroadcast
	MessageTypeGroup     MessageType = entity.RecordTypeGroup
	MessageTypePrivate   MessageType = entity.RecordTypePrivate
)

type Message struct {
	Type     MessageType
	Metadata string

	From    string
	To      string
	Content string
}

type Hub interface {
	RegisterClient(ctx context.Context, k string, conn *websocket.Conn) error
	UnregisterClient(ctx context.Context, k string) error
	SendMessage(ctx context.Context, m Message) error
	Close() error
}

type Client struct {
	conn    *websocket.Conn
	loginAt time.Time
}

func NewHub(env environment.Env, logger logger.Logger, record record.Record) (Hub, error) {
	return &hub{
		clients: make(map[string]*Client),
		record:  record,
	}, nil
}

type hub struct {
	clients map[string]*Client
	record  record.Record
}

func (h *hub) RegisterClient(ctx context.Context, k string, conn *websocket.Conn) error {
	client, ok := h.clients[k]
	if ok {
		client.conn.Close()
		delete(h.clients, k)
	}

	client = &Client{
		conn:    conn,
		loginAt: time.Now(),
	}

	h.clients[k] = client
	return nil
}

func (h *hub) UnregisterClient(ctx context.Context, k string) error {
	client, ok := h.clients[k]
	if ok {
		client.conn.Close()
		delete(h.clients, k)
	}

	return nil
}

func (h *hub) SendBroadcastMessage(ctx context.Context, from, content string) error {
	return nil

}

func (h *hub) SendGroupMessage(ctx context.Context, from, content string, groupID int64) error {
	return nil
}

func (h *hub) SendPrivateMessage(ctx context.Context, from, content, to string) error {
	return nil
}

func (h *hub) SendMessage(ctx context.Context, m Message) error {
	// TODO: add message to storage
	client, ok := h.clients[m.To]
	if !ok {
		// TODO: 用户当前不在线
		return nil
	}

	if err := client.conn.WriteJSON(m); err != nil {
		// TODO: into message queue
		return err
	}

	return nil
}

func (h *hub) Close() error {
	for _, c := range h.clients {
		c.conn.Close()
	}

	return nil
}
