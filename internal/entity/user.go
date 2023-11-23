package entity

import "time"

type User struct {
	Subject  string
	Nickname string
	Username string
	Password string
	Phone    string

	CreatedAt time.Time
}
