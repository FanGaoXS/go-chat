package entity

import (
	"database/sql/driver"
	"time"
)

type FriendRequestLog struct {
	ID       int64                  `json:"id"`
	Sender   string                 `json:"sender"`
	Receiver string                 `json:"receiver"`
	Status   FriendRequestLogStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
}

type FriendRequestLogStatus int

const (
	FriendRequestLogStatusPending FriendRequestLogStatus = iota
	FriendRequestLogStatusAgreed
	FriendRequestLogStatusRefused
)

var friendRequestLogStatusString = map[FriendRequestLogStatus]string{
	FriendRequestLogStatusPending: "待处理",
	FriendRequestLogStatusAgreed:  "已接受",
	FriendRequestLogStatusRefused: "已拒绝",
}

var friendRequestLogStatusID = map[string]FriendRequestLogStatus{
	"待处理": FriendRequestLogStatusPending,
	"已接受": FriendRequestLogStatusAgreed,
	"已拒绝": FriendRequestLogStatusRefused,
}

func (f FriendRequestLogStatus) String() string {
	return friendRequestLogStatusString[f]
}

func (f FriendRequestLogStatus) Value() (driver.Value, error) {
	return friendRequestLogStatusString[f], nil
}

func (f *FriendRequestLogStatus) Scan(value interface{}) error {
	*f = friendRequestLogStatusID[value.(string)]
	return nil
}

func FriendRequestLogStatusFromString(s string) (FriendRequestLogStatus, bool) {
	v, ok := friendRequestLogStatusID[s]
	return v, ok
}
