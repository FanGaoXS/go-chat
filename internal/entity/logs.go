package entity

import (
	"database/sql/driver"
	"time"
)

type LogsStatus int

const (
	LogsStatusPending LogsStatus = iota
	LogsStatusAgreed
	LogsStatusRefused
)

var logsStatusString = map[LogsStatus]string{
	LogsStatusPending: "待处理",
	LogsStatusAgreed:  "已接受",
	LogsStatusRefused: "已拒绝",
}

var logsStatusID = map[string]LogsStatus{
	"待处理": LogsStatusPending,
	"已接受": LogsStatusAgreed,
	"已拒绝": LogsStatusRefused,
}

func (f LogsStatus) String() string {
	return logsStatusString[f]
}

func (f LogsStatus) Value() (driver.Value, error) {
	return logsStatusString[f], nil
}

func (f *LogsStatus) Scan(value interface{}) error {
	*f = logsStatusID[value.(string)]
	return nil
}

func LogsStatusFromString(s string) (LogsStatus, bool) {
	v, ok := logsStatusID[s]
	return v, ok
}

type FriendRequestLog struct {
	ID       int64      `json:"id"`
	Sender   string     `json:"sender"`
	Receiver string     `json:"receiver"`
	Status   LogsStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
}

type GroupInvitationLog struct {
	ID       int64      `json:"id"`
	GroupID  int64      `json:"group_id"`
	Sender   string     `json:"sender"`
	Receiver string     `json:"receiver"`
	Status   LogsStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
}

type GroupRequestLog struct {
	ID       int64      `json:"id"`
	GroupID  int64      `json:"group_id"`
	Sender   string     `json:"sender"`
	Approver NullString `json:"approver"`
	Status   LogsStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
}
