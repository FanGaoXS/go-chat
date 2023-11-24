package entity

import "time"

type User struct {
	Subject  string `json:"subject"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`

	CreatedAt time.Time `json:"created_at"`
}
