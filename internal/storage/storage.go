package storage

import (
	"context"

	"fangaoxs.com/go-chat/internal/entity"
)

type Storage interface {
	Close() error
	NewSession(ctx context.Context) (Session, error)

	InsertUser(ses Session, i *entity.User) error
	GetUserBySubject(ses Session, subject string) (*entity.User, error)
	GetUserBySecret(ses Session, username, password string) (*entity.User, error)
	DeleteUser(ses Session, subject string) error

	InsertUserFriend(ses Session, i *entity.UserFriend) error
	ListUserFriendsByUserSubject(ses Session, userSubject string) ([]*entity.UserFriend, error)
	DeleteUserFriend(ses Session, userSubject, friendSubject string) error

	InsertGroup(ses Session, i *entity.Group) (int64, error)
	GetGroupByID(ses Session, id int64) (*entity.Group, error)
	ListGroupsByCreatedBy(ses Session, createdBy string) ([]*entity.Group, error)
	DeleteGroup(ses Session, id int64) error
	UpdateGroupIsPublic(ses Session, id int64, isPublic bool) error

	InsertGroupMember(ses Session, userSubject string, groupID int64) error
	GetGroupMember(ses Session, userSubject string, groupID int64) (*entity.GroupMember, error)
	ListGroupMembersByGroupID(ses Session, groupID int64) ([]*entity.GroupMember, error)
	ListGroupMembersByUserSubject(ses Session, userSubject string) ([]*entity.GroupMember, error)
}
