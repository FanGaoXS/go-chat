package postgres

import (
	"fmt"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/storage"
)

func (p *postgres) InsertFriendRequestLog(ses storage.Session, i *entity.FriendRequestLog) error {
	sqlstr := rebind(`INSERT INTO "friend_request_log" 
                  (sender, receiver, status)
                  VALUES
                  (?, ?, ?);`)
	args := []any{
		i.Sender,
		i.Receiver,
		i.Status,
	}

	var err error
	_, err = ses.Exec(sqlstr, args...)
	if err != nil {
		return wrapPGErrorf(err, "failed to insert friend request log")
	}

	return nil
}

func (p *postgres) listFriendRequestLogs(ses storage.Session, where *entity.Where) ([]*entity.FriendRequestLog, error) {
	projection := []string{
		"id",
		"sender",
		"receiver",
		"status",
		"created_at",
	}

	var args []any
	sqlstr := fmt.Sprintf(`SELECT %s FROM "friend_request_log"`, strings.Join(projection, ", "))
	if where != nil {
		sel, selArgs, err := where.Parse()
		if err != nil {
			return nil, err
		}
		args = append(args, selArgs...)
		sqlstr += sel
	}

	sqlstr = rebind(sqlstr)
	rows, err := ses.Query(sqlstr, args...)
	if err != nil {
		return nil, wrapPGErrorf(err, "failed to list friend requests")
	}
	defer rows.Close()

	var res []*entity.FriendRequestLog
	for rows.Next() {
		r := entity.FriendRequestLog{}
		if err = rows.Scan(&r.ID, &r.Sender, &r.Receiver, &r.Status, &r.CreatedAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan friend request")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) ListFriendRequestLogsBySender(ses storage.Session, sender string) ([]*entity.FriendRequestLog, error) {
	w := &entity.Where{
		FieldNames:  []string{"sender"},
		FieldValues: []any{sender},
	}

	res, err := p.listFriendRequestLogs(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list friend request logs with sender: %s failed", sender)
	}

	return res, nil
}

func (p *postgres) ListFriendRequestLogsByReceiver(ses storage.Session, receiver string) ([]*entity.FriendRequestLog, error) {
	w := &entity.Where{
		FieldNames:  []string{"receiver"},
		FieldValues: []any{receiver},
	}

	res, err := p.listFriendRequestLogs(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list friend request logs with receiver: %s failed", receiver)
	}

	return res, nil
}

func (p *postgres) GetPendingFriendRequestLog(ses storage.Session, sender, receiver string) (*entity.FriendRequestLog, error) {
	w := &entity.Where{
		FieldNames:  []string{"sender", "receiver", "status"},
		FieldValues: []any{sender, receiver, entity.FriendRequestLogStatusPending},
	}

	res, err := p.listFriendRequestLogs(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "get pending friend request log with sender: %s and receiver: %s failed", sender, receiver)
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no pending friend request log with sender: %s and receiver: %s found", sender, receiver)
	}

	return res[0], nil
}

func (p *postgres) GetFriendRequestLogByID(ses storage.Session, id int64) (*entity.FriendRequestLog, error) {
	w := &entity.Where{
		FieldNames:  []string{"id"},
		FieldValues: []any{id},
	}

	res, err := p.listFriendRequestLogs(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "get friend request log with id: %d failed", id)
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no friend request log with id: %d found", id)
	}

	return res[0], nil
}

func (p *postgres) GetFriendRequestLogByIDForUpdate(ses storage.Session, id int64) (*entity.FriendRequestLog, error) {
	sqlstr := rebind(`SELECT * FROM "friend_request_log" WHERE id = ? FOR UPDATE;`)

	var res entity.FriendRequestLog
	var err error
	err = ses.QueryRow(sqlstr, id).Scan(
		&id, &res.Sender, &res.Receiver, &res.Status, &res.CreatedAt,
	)
	if err != nil {
		return nil, wrapPGErrorf(err, "get friend request log for update with id: %d failed", id)
	}

	return &res, nil
}

func (p *postgres) GetPendingFriendRequestLogForUpdate(ses storage.Session, sender, receiver string) (*entity.FriendRequestLog, error) {
	sqlstr := rebind(`SELECT * 
                            FROM "friend_request_log" 
                            WHERE sender = ? 
                            AND receiver = ? 
                            AND status = ? 
                            FOR UPDATE;`)

	var res entity.FriendRequestLog
	var err error
	err = ses.QueryRow(sqlstr, sender, receiver, entity.FriendRequestLogStatusPending).Scan(
		&res.ID, &res.Sender, &res.Receiver, &res.Status, &res.CreatedAt,
	)
	if err != nil {
		return nil, wrapPGErrorf(err, "get pending friend request log for update with sender: %s and receiver: %s  failed", sender, receiver)
	}

	return &res, nil
}

func (p *postgres) UpdateFriendRequestLogStatus(ses storage.Session, id int64, status entity.FriendRequestLogStatus) error {
	sqlstr := rebind(`UPDATE "friend_request_log" SET status = ? WHERE id = ?;`)
	_, err := ses.Exec(sqlstr, status, id)
	if err != nil {
		return wrapPGErrorf(err, "update friend request log status with id: %d to %s failed", id, status.String())
	}

	return nil
}
