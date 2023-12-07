package storage

import (
	"context"

	"fangaoxs.com/go-chat/internal/entity"
)

type Storage interface {
	Close() error
	NewSession(ctx context.Context) (Session, error)

	InsertUser(ses Session, i *entity.User) error
	ListAllUsers(ses Session) ([]*entity.User, error)
	GetUserBySubject(ses Session, subject string) (*entity.User, error)
	GetUserBySecret(ses Session, username, password string) (*entity.User, error)
	DeleteUser(ses Session, subject string) error

	InsertUserFriend(ses Session, i *entity.UserFriend) error
	IsFriendOfUser(ses Session, userSubject, friendSubject string) (bool, error)
	ListUserFriendsByUserSubject(ses Session, userSubject string) ([]*entity.UserFriend, error)
	DeleteUserFriend(ses Session, userSubject, friendSubject string) error

	InsertGroup(ses Session, i *entity.Group) (int64, error)
	GetGroupByID(ses Session, id int64) (*entity.Group, error)
	ListGroupsByCreatedBy(ses Session, createdBy string) ([]*entity.Group, error)
	DeleteGroup(ses Session, id int64) error
	UpdateGroupIsPublic(ses Session, id int64, isPublic bool) error

	InsertGroupMember(ses Session, i *entity.GroupMember) error
	DeleteGroupMembersByGroupID(ses Session, groupID int64) error
	DeleteGroupMember(ses Session, userSubject string, groupID int64) error
	GetGroupMember(ses Session, userSubject string, groupID int64) (*entity.GroupMember, error)
	IsMemberOfGroup(ses Session, userSubject string, groupID int64) (bool, error)
	ListGroupMembersByGroupID(ses Session, groupID int64) ([]*entity.GroupMember, error)
	ListGroupMembersByUserSubject(ses Session, userSubject string) ([]*entity.GroupMember, error)
	IsAdminOfGroup(ses Session, subject string, groupID int64) (bool, error)
	UpdateGroupMemberIsAdmin(ses Session, subject string, groupID int64, isAdmin bool) error
	ListGroupAdminsByGroupID(ses Session, groupID int64) ([]*entity.GroupMember, error)

	InsertRecordBroadcast(ses Session, i *entity.RecordBroadcast) (int64, error)
	ListAllRecordBroadcasts(ses Session) ([]*entity.RecordBroadcast, error)
	ListRecordBroadcastsBySender(ses Session, sender string) ([]*entity.RecordBroadcast, error)

	InsertRecordGroup(ses Session, i *entity.RecordGroup) (int64, error)
	ListRecordGroupsByGroup(ses Session, groupID int64) ([]*entity.RecordGroup, error)

	InsertRecordPrivate(ses Session, i *entity.RecordPrivate) (int64, error)
	ListRecordPrivatesByParty(ses Session, subject1, subject2 string) ([]*entity.RecordPrivate, error)
}
