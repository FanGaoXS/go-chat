package entity

import (
	"database/sql/driver"
	"time"
)

type Record struct {
	SeqID    string     `json:"seq_id"`
	Type     RecordType `json:"type"`
	Metadata string     `json:"metadata"`

	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type RecordType int

const (
	RecordTypeInvalid RecordType = iota
	RecordTypeBroadcast
	RecordTypeGroup
	RecordTypePrivate
)

var recordTypeString = map[RecordType]string{
	RecordTypeInvalid:   "非法",
	RecordTypeBroadcast: "广播",
	RecordTypeGroup:     "群聊",
	RecordTypePrivate:   "私聊",
}

var recordTypeID = map[string]RecordType{
	"非法": RecordTypeInvalid,
	"广播": RecordTypeBroadcast,
	"群聊": RecordTypeGroup,
	"私聊": RecordTypePrivate,
}

func (r RecordType) String() string {
	return recordTypeString[r]
}

func (r RecordType) Value() (driver.Value, error) {
	return recordTypeString[r], nil
}

func (r *RecordType) Scan(value interface{}) error {
	*r = recordTypeID[value.(string)]
	return nil
}

func RecordTypeFromString(s string) (RecordType, bool) {
	v, ok := recordTypeID[s]
	return v, ok
}
