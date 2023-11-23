package entity

import (
	"database/sql/driver"
	"time"
)

type Group struct {
	ID        int64
	Name      string
	Type      GroupType
	IsPublic  bool // 是否公开的群
	CreatedBy string
	CreatedAt time.Time
}

type GroupType int

const (
	DefaultGroupType GroupType = iota
	DatingGroupType
	GameGroupType
	StudyGroupType
)

var groupTypeString = map[GroupType]string{
	DefaultGroupType: "默认",
	DatingGroupType:  "交友",
	GameGroupType:    "游戏",
	StudyGroupType:   "学习",
}

var groupTypeID = map[string]GroupType{
	"默认": DefaultGroupType,
	"交友": DatingGroupType,
	"游戏": GameGroupType,
	"学习": StudyGroupType,
}

func (g GroupType) String() string {
	return groupTypeString[g]
}

func GroupTypeFromString(s string) (GroupType, bool) {
	v, ok := groupTypeID[s]
	return v, ok
}

func (g GroupType) Value() (driver.Value, error) {
	return groupTypeString[g], nil
}

func (g *GroupType) Scan(value interface{}) error {
	*g = groupTypeID[value.(string)]
	return nil
}
