package entity

import (
	"database/sql/driver"
	"time"
)

type User struct {
	Subject  string `json:"subject"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Password string `json:"-"`
	Phone    string `json:"phone"`

	CreatedAt time.Time `json:"created_at"`
}

type Friendship struct {
	UserSubject   string `json:"user_subject"`
	FriendSubject string `json:"friend_subject"`

	CreatedAt time.Time `json:"created_at"`
}

type FriendRequest struct {
	Sender   string              `json:"sender"`
	Receiver string              `json:"receiver"`
	Status   FriendRequestStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
}

type FriendRequestStatus int

const (
	FriendRequestStatusPending FriendRequestStatus = iota
	FriendRequestStatusAgreed
	FriendRequestStatusRefused
)

var friendRequestStatusString = map[FriendRequestStatus]string{
	FriendRequestStatusPending: "待处理",
	FriendRequestStatusAgreed:  "已接受",
	FriendRequestStatusRefused: "已拒绝",
}

var friendRequestStatusID = map[string]FriendRequestStatus{
	"待处理": FriendRequestStatusPending,
	"已接受": FriendRequestStatusAgreed,
	"已拒绝": FriendRequestStatusRefused,
}

func (f FriendRequestStatus) String() string {
	return friendRequestStatusString[f]
}

func (f FriendRequestStatus) Value() (driver.Value, error) {
	return friendRequestStatusString[f], nil
}

func (f *FriendRequestStatus) Scan(value interface{}) error {
	*f = friendRequestStatusID[value.(string)]
	return nil
}

func FriendRequestStatusFromString(s string) (FriendRequestStatus, bool) {
	v, ok := friendRequestStatusID[s]
	return v, ok
}
