package entity

import "time"

type RecordBroadcast struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
	Sender  string `json:"sender"`

	CreatedAt time.Time `json:"created_at"`
}

type RecordGroup struct {
	ID      int64  `json:"id"`
	GroupID int64  `json:"group_id"`
	Content string `json:"content"`
	Sender  string `json:"sender"`

	CreatedAt time.Time `json:"created_at"`
}

type RecordPrivate struct {
	ID       int64  `json:"id"`
	Content  string `json:"content"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`

	CreatedAt time.Time `json:"created_at"`
}
