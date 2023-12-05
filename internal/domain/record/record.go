package record

import (
	"context"
	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"fangaoxs.com/go-chat/internal/storage"
	"strings"
)

type Record interface {
	InsertRecordBroadcast(ctx context.Context, sender, content string) error
	InsertRecordGroup(ctx context.Context, sender, content string, groupID int64) error
	InsertRecordPrivate(ctx context.Context, sender, content, receiver string) error

	ListRecordBroadcasts(ctx context.Context, sender string) ([]*entity.RecordBroadcast, error)
	ListRecordBroadcastsBySender(ctx context.Context, sender string) ([]*entity.RecordBroadcast, error)
	ListRecordGroups(ctx context.Context, groupID int64) ([]*entity.RecordGroup, error)
	ListRecordPrivate(ctx context.Context, sender, receiver string) ([]*entity.RecordPrivate, error)
}

func New(env environment.Env, logger logger.Logger, storage storage.Storage) (Record, error) {
	return &record{
		logger:  logger,
		storage: storage,
	}, nil
}

type record struct {
	logger logger.Logger

	storage storage.Storage
}

func (r *record) InsertRecordBroadcast(ctx context.Context, sender, content string) error {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return err
	}

	rcd := &entity.RecordBroadcast{
		Content: content,
		Sender:  sender,
	}
	_, err = r.storage.InsertRecordBroadcast(ses, rcd)
	if err != nil {
		return err
	}

	return nil
}

func (r *record) InsertRecordGroup(ctx context.Context, sender, content string, groupID int64) error {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return err
	}

	// 检查sender是否在group中
	_, err = r.storage.GetGroupMember(ses, sender, groupID)
	if err != nil {
		return err
	}

	rcd := &entity.RecordGroup{
		GroupID: groupID,
		Content: content,
		Sender:  sender,
	}
	_, err = r.storage.InsertRecordGroup(ses, rcd)
	if err != nil {
		return err
	}

	return nil
}

func (r *record) InsertRecordPrivate(ctx context.Context, sender, content, receiver string) error {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return err
	}

	// 检查receiver是否存在
	_, err = r.storage.GetUserBySubject(ses, receiver)
	if err != nil {
		return err
	}

	ok, err := r.storage.IsFriendOfUser(ses, sender, receiver)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Newf(errors.NotFound, nil, "发送者没有添加接受者为好友")
	}

	ok, err = r.storage.IsFriendOfUser(ses, receiver, sender)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Newf(errors.NotFound, nil, "接受者没有添加发送者为好友")
	}

	rcd := &entity.RecordPrivate{
		Content:  content,
		Sender:   sender,
		Receiver: receiver,
	}
	_, err = r.storage.InsertRecordPrivate(ses, rcd)
	if err != nil {
		return err
	}

	return nil
}

func (r *record) ListRecordBroadcasts(ctx context.Context, sender string) ([]*entity.RecordBroadcast, error) {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	var res []*entity.RecordBroadcast
	if sender = strings.TrimSpace(sender); sender == "" {
		res, err = r.storage.ListAllRecordBroadcasts(ses)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = r.storage.ListRecordBroadcastsBySender(ses, sender)
		if err != nil {
			return nil, err
		}
	}
	if len(res) == 0 {
		return nil, errors.New(errors.NotFound, nil, "empty record_broadcast")
	}

	return res, nil
}

func (r *record) ListRecordBroadcastsBySender(ctx context.Context, sender string) ([]*entity.RecordBroadcast, error) {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := r.storage.ListRecordBroadcastsBySender(ses, sender)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "empty record_broadcast with sender: %s", sender)
	}

	return res, nil
}

func (r *record) ListRecordGroups(ctx context.Context, groupID int64) ([]*entity.RecordGroup, error) {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := r.storage.ListRecordGroupsByGroup(ses, groupID)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "empty record_group with group: %d", groupID)
	}

	return res, nil
}

func (r *record) ListRecordPrivate(ctx context.Context, subject1, subject2 string) ([]*entity.RecordPrivate, error) {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := r.storage.ListRecordPrivatesByParty(ses, subject1, subject2)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "empty record_private with subject1: %s and subject2: %s", subject1, subject2)
	}

	return res, nil
}
