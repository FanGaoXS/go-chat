package hub

import (
	"context"
	"time"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/record"
	"fangaoxs.com/go-chat/internal/infras/logger"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	loginAt time.Time
}

type Hub interface {
	Close() error

	RegisterClient(ctx context.Context, subject string, conn *websocket.Conn) error
	UnregisterClient(ctx context.Context, subject string) error

	SendBroadcastMessage(ctx context.Context, sender, content string) error
	SendGroupMessage(ctx context.Context, sender, content string, groupID int64) error
	SendPrivateMessage(ctx context.Context, sender, content, receiver string) error
}

func NewHub(env environment.Env, logger logger.Logger, record record.Record, group group.Group) (Hub, error) {
	return &hub{
		clients: make(map[string]*Client),
		record:  record,
		group:   group,
	}, nil
}

type hub struct {
	clients map[string]*Client

	record record.Record
	group  group.Group
}

func (h *hub) Close() error {
	for _, c := range h.clients {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "服务器关闭"))
		c.conn.Close()
	}

	return nil
}

func (h *hub) RegisterClient(ctx context.Context, subject string, conn *websocket.Conn) error {
	c, ok := h.clients[subject]
	if ok {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "你被强制下线"))
		c.conn.Close()
		delete(h.clients, subject)
	}

	c = &Client{
		conn:    conn,
		loginAt: time.Now(),
	}

	h.clients[subject] = c
	return nil
}

func (h *hub) UnregisterClient(ctx context.Context, subject string) error {
	c, ok := h.clients[subject]
	if ok {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "注销"))
		c.conn.Close()
		delete(h.clients, subject)
	}

	return nil
}

func (h *hub) SendBroadcastMessage(ctx context.Context, sender, content string) error {
	if err := h.record.InsertRecordBroadcast(ctx, sender, content); err != nil {
		return err
	}

	for _, c := range h.clients {
		m := map[string]any{
			"type":    "broadcast",
			"content": content,
			"sender":  sender,
		}
		err := c.conn.WriteJSON(m)
		if err != nil {
			// TODO: 重试
		}
	}

	return nil

}

func (h *hub) SendGroupMessage(ctx context.Context, sender, content string, groupID int64) error {
	if err := h.record.InsertRecordGroup(ctx, sender, content, groupID); err != nil {
		return err
	}

	members, err := h.group.ListMembersOfGroup(ctx, groupID)
	if err != nil {
		return err
	}

	for _, member := range members {
		if member.Subject == sender {
			// 不发送给自己
			continue
		}
		c, ok := h.clients[member.Subject]
		if !ok {
			// 群成员不在线
			continue
		}
		m := map[string]any{
			"type":     "group",
			"group_id": groupID,
			"content":  content,
			"sender":   sender,
		}
		err = c.conn.WriteJSON(m)
		if err != nil {
			// TODO: 重试
		}
	}

	return nil
}

func (h *hub) SendPrivateMessage(ctx context.Context, sender, content, receiver string) error {
	if err := h.record.InsertRecordPrivate(ctx, sender, content, receiver); err != nil {
		return err
	}

	c, ok := h.clients[receiver]
	if !ok {
		// 对方不在线
		return nil
	}
	m := map[string]any{
		"type":     "private",
		"content":  content,
		"sender":   sender,
		"receiver": receiver,
	}
	err := c.conn.WriteJSON(m)
	if err != nil {
		// TODO: 重试
	}

	return nil
}
