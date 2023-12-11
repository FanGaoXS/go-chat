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
	CreateFriendApplication(ctx context.Context, sender, receiver string) error
	AgreeFriendApplication(ctx context.Context, id int64) error
	RefuseFriendApplication(ctx context.Context, id int64) error

	GetFriendApplication(ctx context.Context, id int64) (*entity.FriendRequestLog, error)
	FriendApplicationsFrom(ctx context.Context, subject string) ([]*entity.FriendRequestLog, error)
	FriendApplicationsTo(ctx context.Context, subject string) ([]*entity.FriendRequestLog, error)
}

func New(env environment.Env, logger logger.Logger, storage storage.Storage) (Applications, error) {
	return &applications{storage: storage}, nil
}

type applications struct {
	storage storage.Storage
}

func (a *applications) CreateFriendApplication(ctx context.Context, sender, receiver string) error {
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
		Status:   entity.FriendRequestLogStatusPending,
	}

	if got != nil {
		if err = a.storage.UpdateFriendRequestLogStatus(ses, got.ID, entity.FriendRequestLogStatusAgreed); err != nil {
			return err
		}

		if err = a.storage.InsertFriendship(ses, sender, receiver); err != nil {
			return err
		}
		if err = a.storage.InsertFriendship(ses, receiver, sender); err != nil {
			return err
		}
		requestlog.Status = entity.FriendRequestLogStatusAgreed
	}

	if err = a.storage.InsertFriendRequestLog(ses, requestlog); err != nil {
		return err
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *applications) AgreeFriendApplication(ctx context.Context, id int64) error {
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
	if requsetlog.Status != entity.FriendRequestLogStatusPending {
		return errors.Newf(errors.InvalidArgument, nil, "好友申请请求已经被处理")
	}
	err = a.storage.UpdateFriendRequestLogStatus(ses, id, entity.FriendRequestLogStatusAgreed)
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

func (a *applications) RefuseFriendApplication(ctx context.Context, id int64) error {
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
	if requsetlog.Status != entity.FriendRequestLogStatusPending {
		return errors.Newf(errors.InvalidArgument, nil, "该好友申请请求已经被处理")
	}
	err = a.storage.UpdateFriendRequestLogStatus(ses, id, entity.FriendRequestLogStatusRefused)
	if err != nil {
		return err
	}

	if err = ses.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *applications) GetFriendApplication(ctx context.Context, id int64) (*entity.FriendRequestLog, error) {
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

func (a *applications) FriendApplicationsFrom(ctx context.Context, subject string) ([]*entity.FriendRequestLog, error) {
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

func (a *applications) FriendApplicationsTo(ctx context.Context, subject string) ([]*entity.FriendRequestLog, error) {
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
