package applications

import (
	"context"
	
	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"fangaoxs.com/go-chat/internal/storage"
)

type Applications interface {
	// FriendRequest 好友申请

	CreateFriendRequest(ctx context.Context, sender, receiver string) error
	AgreeFriendRequest(ctx context.Context, id int64, approver string) error
	RefuseFriendRequest(ctx context.Context, id int64, approver string) error
	GetFriendRequest(ctx context.Context, id int64) (*entity.FriendRequestLog, error)
	FriendRequestsFrom(ctx context.Context, subject string) ([]*entity.FriendRequestLog, error)
	FriendRequestsTo(ctx context.Context, subject string) ([]*entity.FriendRequestLog, error)

	// GroupInvitation 邀请入群

	CreateGroupInvitation(ctx context.Context, sender, receiver string, groupID int64) error
	AgreeGroupInvitation(ctx context.Context, id int64) error
	RefuseGroupInvitation(ctx context.Context, id int64) error
	GetGroupInvitation(ctx context.Context, id int64) (*entity.GroupInvitationLog, error)
	GroupInvitationsTo(ctx context.Context, receiver string) ([]*entity.GroupInvitationLog, error)

	// GroupRequest 申请加群

	CreateGroupRequest(ctx context.Context, sender string, groupID int64) error
	AgreeGroupRequest(ctx context.Context, id int64, approver string) error
	RefuseGroupRequest(ctx context.Context, id int64, approver string) error
	GetGroupRequest(ctx context.Context, id int64) (*entity.GroupRequestLog, error)
	GroupRequestsFrom(ctx context.Context, sender string) ([]*entity.GroupRequestLog, error)
	GroupRequestsTo(ctx context.Context, groupID int64) ([]*entity.GroupRequestLog, error)
}

func New(env environment.Env, logger logger.Logger, storage storage.Storage) (Applications, error) {
	return &applications{storage: storage}, nil
}

type applications struct {
	storage storage.Storage
}

func (a *applications) CreateFriendRequest(ctx context.Context, sender, receiver string) error {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	if sender == receiver {
		return errors.Newf(errors.InvalidArgument, nil, "不可以向自己发起好友申请")
	}

	_, err = a.storage.GetUserBySubject(ses, receiver)
	if err != nil {
		return err
	}

	ok, err := a.storage.IsFriendOfUser(ses, sender, receiver)
	if err != nil {
		return err
	}
	if ok {
		return errors.Newf(errors.AlreadyExists, nil, "[%s]已经是[%s]的好友", receiver, sender)
	}

	got, err := a.storage.GetPendingFriendRequestLog(ses, sender, receiver)
	if err != nil && errors.Code(err) != errors.NotFound {
		return err
	}
	if got != nil {
		return errors.Newf(errors.AlreadyExists, nil, "已经存在[%s]发送给[%s]的未处理的好友请求", sender, receiver)
	}

	// 不存在sender发送给receiver的好友请求

	// 检查是否已经有receiver发送给sender的且pending的好友请求，如果有，则直接双向同意
	got, err = a.storage.GetPendingFriendRequestLogForUpdate(ses, receiver, sender)
	if err != nil && errors.Code(err) != errors.NotFound {
		return err
	}

	requestlog := &entity.FriendRequestLog{
		Sender:   sender,
		Receiver: receiver,
		Status:   entity.LogsStatusPending,
	}

	if got != nil {
		if err = a.storage.UpdateFriendRequestLogStatus(ses, got.ID, entity.LogsStatusAgreed); err != nil {
			return err
		}

		if err = a.storage.InsertFriendship(ses, sender, receiver); err != nil {
			return err
		}
		if err = a.storage.InsertFriendship(ses, receiver, sender); err != nil {
			return err
		}
		requestlog.Status = entity.LogsStatusAgreed
	}

	if err = a.storage.InsertFriendRequestLog(ses, requestlog); err != nil {
		return err
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *applications) AgreeFriendRequest(ctx context.Context, id int64, approver string) error {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	requsetlog, err := a.storage.GetFriendRequestLogByIDForUpdate(ses, id)
	if err != nil {
		return err
	}
	if requsetlog.Status != entity.LogsStatusPending {
		return errors.Newf(errors.InvalidArgument, nil, "好友申请请求已经被处理")
	}
	err = a.storage.UpdateFriendRequestLogStatus(ses, id, entity.LogsStatusAgreed)
	if err != nil {
		return err
	}

	if err = a.storage.InsertFriendship(ses, requsetlog.Sender, requsetlog.Receiver); err != nil {
		return err
	}
	if err = a.storage.InsertFriendship(ses, requsetlog.Receiver, requsetlog.Sender); err != nil {
		return err
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *applications) RefuseFriendRequest(ctx context.Context, id int64, approver string) error {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	requsetlog, err := a.storage.GetFriendRequestLogByIDForUpdate(ses, id)
	if err != nil {
		return err
	}
	if requsetlog.Status != entity.LogsStatusPending {
		return errors.Newf(errors.InvalidArgument, nil, "该好友申请请求已经被处理")
	}
	err = a.storage.UpdateFriendRequestLogStatus(ses, id, entity.LogsStatusRefused)
	if err != nil {
		return err
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *applications) GetFriendRequest(ctx context.Context, id int64) (*entity.FriendRequestLog, error) {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := a.storage.GetFriendRequestLogByID(ses, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *applications) FriendRequestsFrom(ctx context.Context, subject string) ([]*entity.FriendRequestLog, error) {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := a.storage.ListFriendRequestLogsBySender(ses, subject)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no friend applications from %s", subject)
	}

	return res, nil
}

func (a *applications) FriendRequestsTo(ctx context.Context, subject string) ([]*entity.FriendRequestLog, error) {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := a.storage.ListFriendRequestLogsByReceiver(ses, subject)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no friend applications to %s", subject)
	}

	return res, nil
}

func (a *applications) CreateGroupInvitation(ctx context.Context, sender, receiver string, groupID int64) error {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	if sender == receiver {
		return errors.Newf(errors.InvalidArgument, nil, "不可以邀请自己入群")
	}

	_, err = a.storage.GetUserBySubject(ses, receiver)
	if err != nil {
		return err
	}

	_, err = a.storage.GetGroupByID(ses, groupID)
	if err != nil {
		return err
	}

	ok, err := a.storage.IsMemberOfGroup(ses, receiver, groupID)
	if err != nil {
		return err
	}
	if ok {
		return errors.Newf(errors.AlreadyExists, nil, "[%s]已经是群[%d]成员了", receiver, groupID)
	}

	got, err := a.storage.GetPendingGroupInvitationLog(ses, groupID, receiver)
	if err != nil && errors.Code(err) != errors.NotFound {
		return err
	}
	if got != nil {
		return errors.Newf(errors.AlreadyExists, nil, "已经存在[%s]发送给[%s]的邀请加入群[%d]的未处理请求", got.Sender, receiver, groupID)
	}

	// 不存在sender发送给receiver的邀请入groupID群的请求

	// 检查是否已经有receiver发送给groupID的且pending的入群，如果有，则直接双向同意
	forupdate, err := a.storage.GetPendingGroupRequestLogForUpdate(ses, groupID, receiver)
	if err != nil && errors.Code(err) != errors.NotFound {
		return err
	}

	insert := &entity.GroupInvitationLog{
		GroupID:  groupID,
		Sender:   sender,
		Receiver: receiver,
		Status:   entity.LogsStatusPending,
	}

	if forupdate != nil {
		if err = a.storage.UpdateGroupRequestLogStatus(ses, forupdate.ID, sender, entity.LogsStatusAgreed); err != nil {
			return err
		}

		e := &entity.GroupMember{
			UserSubject: receiver,
			GroupID:     groupID,
			IsAdmin:     false,
		}
		if err = a.storage.InsertGroupMember(ses, e); err != nil {
			return err
		}
		insert.Status = entity.LogsStatusAgreed
	}

	if err = a.storage.InsertGroupInvitationLog(ses, insert); err != nil {
		return err
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *applications) AgreeGroupInvitation(ctx context.Context, id int64) error {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	forupdate, err := a.storage.GetGroupInvitationLogByIDForUpdate(ses, id)
	if err != nil {
		return err
	}
	if forupdate.Status != entity.LogsStatusPending {
		return errors.Newf(errors.InvalidArgument, nil, "邀请入群已经被处理")
	}
	if err = a.storage.UpdateGroupInvitationLogStatus(ses, id, entity.LogsStatusAgreed); err != nil {
		return err
	}

	i := &entity.GroupMember{
		UserSubject: forupdate.Receiver,
		GroupID:     forupdate.GroupID,
		IsAdmin:     false,
	}
	if err = a.storage.InsertGroupMember(ses, i); err != nil {
		return err
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *applications) RefuseGroupInvitation(ctx context.Context, id int64) error {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	forupdate, err := a.storage.GetGroupInvitationLogByIDForUpdate(ses, id)
	if err != nil {
		return err
	}
	if forupdate.Status != entity.LogsStatusPending {
		return errors.Newf(errors.InvalidArgument, nil, "邀请入群已经被处理")
	}
	err = a.storage.UpdateGroupInvitationLogStatus(ses, id, entity.LogsStatusRefused)
	if err != nil {
		return err
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *applications) GetGroupInvitation(ctx context.Context, id int64) (*entity.GroupInvitationLog, error) {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := a.storage.GetGroupInvitationLog(ses, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *applications) GroupInvitationsFrom(ctx context.Context, sender string) ([]*entity.GroupInvitationLog, error) {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := a.storage.ListGroupInvitationLogsBySender(ses, sender)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no invitations from [%s] found", sender)
	}

	return res, nil
}

func (a *applications) GroupInvitationsTo(ctx context.Context, receiver string) ([]*entity.GroupInvitationLog, error) {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := a.storage.ListGroupInvitationLogsByReceiver(ses, receiver)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no invitations to [%s] found", receiver)
	}

	return res, nil
}

func (a *applications) CreateGroupRequest(ctx context.Context, sender string, groupID int64) error {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	g, err := a.storage.GetGroupByID(ses, groupID)
	if err != nil {
		return err
	}

	if !g.IsPublic {
		return errors.Newf(errors.PermissionDenied, nil, "群[%d]是非公开的群组", groupID)
	}

	ok, err := a.storage.IsMemberOfGroup(ses, sender, groupID)
	if err != nil {
		return err
	}
	if ok {
		return errors.Newf(errors.AlreadyExists, nil, "[%s]已经是群[%d]成员了", sender, groupID)
	}

	got, err := a.storage.GetPendingGroupRequestLog(ses, groupID, sender)
	if err != nil && errors.Code(err) != errors.NotFound {
		return err
	}
	if got != nil {
		return errors.Newf(errors.AlreadyExists, nil, "已经存在[%s]发送给群[%d]的申请入群请求了", sender, groupID)
	}

	// 不存在sender发送给groupID群的入群请求

	// 检查是否已经有group邀请sender的且pending的入群请求，如果有，则双向同意
	forUpdate, err := a.storage.GetPendingGroupInvitationLogForUpdate(ses, groupID, sender)
	if err != nil && errors.Code(err) != errors.NotFound {
		return err
	}

	insert := &entity.GroupRequestLog{
		GroupID: groupID,
		Sender:  sender,
		Status:  entity.LogsStatusPending,
	}

	if forUpdate != nil {
		if err = a.storage.UpdateGroupInvitationLogStatus(ses, forUpdate.ID, entity.LogsStatusAgreed); err != nil {
			return err
		}

		e := &entity.GroupMember{
			UserSubject: sender,
			GroupID:     groupID,
			IsAdmin:     false,
		}
		if err = a.storage.InsertGroupMember(ses, e); err != nil {
			return err
		}
		insert.Status = entity.LogsStatusAgreed
	}

	if err = a.storage.InsertGroupRequestLog(ses, insert); err != nil {
		return err
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *applications) AgreeGroupRequest(ctx context.Context, id int64, approver string) error {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	forUpdate, err := a.storage.GetGroupRequestLogByIDForUpdate(ses, id)
	if err != nil {
		return err
	}
	if forUpdate.Status != entity.LogsStatusPending {
		return errors.Newf(errors.InvalidArgument, nil, "入群申请已经被处理")
	}
	if err = a.storage.UpdateGroupRequestLogStatus(ses, id, approver, entity.LogsStatusAgreed); err != nil {
		return err
	}

	i := &entity.GroupMember{
		UserSubject: forUpdate.Sender,
		GroupID:     forUpdate.GroupID,
		IsAdmin:     false,
	}
	if err = a.storage.InsertGroupMember(ses, i); err != nil {
		return err
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *applications) RefuseGroupRequest(ctx context.Context, id int64, approver string) error {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return err
	}
	ses, err = ses.Begin()
	if err != nil {
		return err
	}

	forUpdate, err := a.storage.GetGroupRequestLogByIDForUpdate(ses, id)
	if err != nil {
		return err
	}
	if forUpdate.Status != entity.LogsStatusPending {
		return errors.Newf(errors.InvalidArgument, nil, "入群申请已经被处理")
	}
	if err = a.storage.UpdateGroupRequestLogStatus(ses, id, approver, entity.LogsStatusRefused); err != nil {
		return err
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *applications) GetGroupRequest(ctx context.Context, id int64) (*entity.GroupRequestLog, error) {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := a.storage.GetGroupRequestLog(ses, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *applications) GroupRequestsFrom(ctx context.Context, sender string) ([]*entity.GroupRequestLog, error) {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := a.storage.ListGroupRequestLogsBySender(ses, sender)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no group requests from %s found", sender)
	}

	return res, nil
}

func (a *applications) GroupRequestsTo(ctx context.Context, groupID int64) ([]*entity.GroupRequestLog, error) {
	ses, err := a.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := a.storage.ListGroupRequestLogsByGroup(ses, groupID)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no group requests to group: %d found", groupID)
	}

	return res, nil
}
