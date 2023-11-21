package entity

import "time"

type User struct {
	ID       int64
	Nickname string
	Username string
	Password string
	Phone    string

	CreatedAt time.Time
}
