package entity

import "time"

type GroupMember struct {
	UserSubject string    `json:"user_subject"`
	GroupID     int64     `json:"group_id"`
	JoinAt      time.Time `json:"join_at"`
}
