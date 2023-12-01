package record

import (
	"context"
	"fmt"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"fangaoxs.com/go-chat/internal/storage"
)

type Record interface {
	ListBroadcastRecords(ctx context.Context) ([]*entity.Record, error)
	ListGroupRecords(ctx context.Context, groupID int64) ([]*entity.Record, error)
	ListPrivateRecord(ctx context.Context, sender, receiver string) ([]*entity.Record, error)
}

func New(storage storage.Storage, logger logger.Logger) (Record, error) {
	return &record{
		storage: storage,
		logger:  logger,
	}, nil
}

type record struct {
	storage storage.Storage
	logger  logger.Logger
}

func (r *record) InsertRecord(ctx context.Context, record *entity.Record) {

}

func (r *record) ListBroadcastRecords(ctx context.Context) ([]*entity.Record, error) {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	records, err := r.storage.ListBroadcastRecords(ses)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no broadcast records found")
	}

	return records, nil
}

func (r *record) ListGroupRecords(ctx context.Context, groupID int64) ([]*entity.Record, error) {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	metadata := fmt.Sprintf("group_id:%d", groupID)
	records, err := r.storage.ListGroupRecords(ses, metadata)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no broadcast records found")
	}

	return records, nil
}

func (r *record) ListPrivateRecord(ctx context.Context, sender, receiver string) ([]*entity.Record, error) {
	ses, err := r.storage.NewSession(ctx)
	if err != nil {
		return nil, err
	}

	records, err := r.storage.ListPrivateRecords(ses, sender, receiver)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no broadcast records found")
	}

	return records, nil
}
