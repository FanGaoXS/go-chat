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

	InsertFriendship(ses Session, userSubject, friendSubject string) error
	IsFriendOfUser(ses Session, userSubject, friendSubject string) (bool, error)
	ListFriendshipsByUserSubject(ses Session, userSubject string) ([]*entity.Friendship, error)
	DeleteFriendship(ses Session, userSubject, friendSubject string) error

	InsertFriendRequestLog(ses Session, i *entity.FriendRequestLog) error
	ListFriendRequestLogsBySender(ses Session, sender string) ([]*entity.FriendRequestLog, error)
	ListFriendRequestLogsByReceiver(ses Session, receiver string) ([]*entity.FriendRequestLog, error)
	GetFriendRequestLogByID(ses Session, id int64) (*entity.FriendRequestLog, error)
	GetPendingFriendRequestLog(ses Session, sender, receiver string) (*entity.FriendRequestLog, error)
	GetFriendRequestLogByIDForUpdate(ses Session, id int64) (*entity.FriendRequestLog, error)
	GetPendingFriendRequestLogForUpdate(ses Session, sender, receiver string) (*entity.FriendRequestLog, error)
	UpdateFriendRequestLogStatus(ses Session, id int64, status entity.LogsStatus) error

	InsertGroupRequestLog(ses Session, i *entity.GroupRequestLog) error
	ListGroupRequestLogsByGroup(ses Session, groupID int64) ([]*entity.GroupRequestLog, error)
	ListGroupRequestLogsBySender(ses Session, sender string) ([]*entity.GroupRequestLog, error)
	GetGroupRequestLog(ses Session, id int64) (*entity.GroupRequestLog, error)
	GetPendingGroupRequestLog(ses Session, groupID int64, sender string) (*entity.GroupRequestLog, error)
	GetPendingGroupRequestLogForUpdate(ses Session, groupID int64, sender string) (*entity.GroupRequestLog, error)
	GetGroupRequestLogByIDForUpdate(ses Session, id int64) (*entity.GroupRequestLog, error)
	UpdateGroupRequestLogStatus(ses Session, id int64, approver string, status entity.LogsStatus) error

	InsertGroupInvitationLog(ses Session, i *entity.GroupInvitationLog) error
	ListGroupInvitationLogsByReceiver(ses Session, receiver string) ([]*entity.GroupInvitationLog, error)
	ListGroupInvitationLogsBySender(ses Session, sender string) ([]*entity.GroupInvitationLog, error)
	GetGroupInvitationLog(ses Session, id int64) (*entity.GroupInvitationLog, error)
	GetPendingGroupInvitationLog(ses Session, groupID int64, receiver string) (*entity.GroupInvitationLog, error)
	GetPendingGroupInvitationLogForUpdate(ses Session, groupID int64, receiver string) (*entity.GroupInvitationLog, error)
	GetGroupInvitationLogByIDForUpdate(ses Session, id int64) (*entity.GroupInvitationLog, error)
	UpdateGroupInvitationLogStatus(ses Session, id int64, status entity.LogsStatus) error

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
