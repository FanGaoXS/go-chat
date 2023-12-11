package records

import (
	"context"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"fangaoxs.com/go-chat/internal/storage"
)

type Records interface {
	InsertRecordBroadcast(ctx context.Context, sender, content string) error
	InsertRecordGroup(ctx context.Context, sender, content string, groupID int64) error
	InsertRecordPrivate(ctx context.Context, sender, content, receiver string) error

	ListAllRecordBroadcasts(ctx context.Context) ([]*entity.RecordBroadcast, error)
	ListRecordBroadcastsBySender(ctx context.Context, sender string) ([]*entity.RecordBroadcast, error)
	ListRecordGroups(ctx context.Context, groupID int64) ([]*entity.RecordGroup, error)
	ListRecordPrivate(ctx context.Context, sender, receiver string) ([]*entity.RecordPrivate, error)
}

func New(env environment.Env, logger logger.Logger, storage storage.Storage) (Records, error) {
	return &records{
		logger:  logger,
		storage: storage,
	}, nil
}

type records struct {
	logger logger.Logger

	storage storage.Storage
}

func (r *records) InsertRecordBroadcast(ctx context.Context, sender, content string) error {
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

func (r *records) InsertRecordGroup(ctx context.Context, sender, content string, groupID int64) error {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return err
	}

	_, err = r.storage.GetGroupByID(ses, groupID)
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

func (r *records) InsertRecordPrivate(ctx context.Context, sender, content, receiver string) error {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return err
	}

	_, err = r.storage.GetUserBySubject(ses, receiver)
	if err != nil {
		return err
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

func (r *records) ListAllRecordBroadcasts(ctx context.Context) ([]*entity.RecordBroadcast, error) {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	res, err := r.storage.ListAllRecordBroadcasts(ses)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New(errors.NotFound, nil, "empty record_broadcast")
	}

	return res, nil
}

func (r *records) ListRecordBroadcastsBySender(ctx context.Context, sender string) ([]*entity.RecordBroadcast, error) {
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

// ListRecordGroups 查询groupID群的群聊记录，当且仅当groupID存在时
func (r *records) ListRecordGroups(ctx context.Context, groupID int64) ([]*entity.RecordGroup, error) {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.storage.GetGroupByID(ses, groupID)
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

// ListRecordPrivate 查询subject1和subject2的私聊记录，当且仅当subject1和subject2存在时
func (r *records) ListRecordPrivate(ctx context.Context, subject1, subject2 string) ([]*entity.RecordPrivate, error) {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.storage.GetUserBySubject(ses, subject1)
	if err != nil {
		return nil, err
	}
	_, err = r.storage.GetUserBySubject(ses, subject2)
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
