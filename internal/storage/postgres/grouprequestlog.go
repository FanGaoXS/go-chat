package postgres

import (
	"fmt"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/storage"
)

func (p *postgres) InsertGroupRequestLog(ses storage.Session, i *entity.GroupRequestLog) error {
	sqlstr := rebind(`INSERT INTO "group_request_log" 
                  (group_id, sender, status)
                  VALUES
                  (?, ?, ?, ?);`)
	args := []any{
		i.GroupID,
		i.Sender,
		i.Status,
	}

	var err error
	_, err = ses.Exec(sqlstr, args...)
	if err != nil {
		return wrapPGErrorf(err, "failed to insert group request log")
	}

	return nil
}

func (p *postgres) listGroupRequestLogs(ses storage.Session, where *entity.Where) ([]*entity.GroupRequestLog, error) {
	projection := []string{
		"id",
		"group_id",
		"sender",
		"status",
		"created_at",
	}

	var args []any
	sqlstr := fmt.Sprintf(`SELECT %s FROM "group_request_log"`, strings.Join(projection, ", "))
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
		return nil, wrapPGErrorf(err, "failed to list group request logs")
	}
	defer rows.Close()

	var res []*entity.GroupRequestLog
	for rows.Next() {
		r := entity.GroupRequestLog{}
		if err = rows.Scan(&r.ID, &r.GroupID, &r.Sender, &r.Status, &r.CreatedAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan group request log")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) ListGroupRequestLogsByGroup(ses storage.Session, groupID int64) ([]*entity.GroupRequestLog, error) {
	w := &entity.Where{
		FieldNames:  []string{"group_id"},
		FieldValues: []any{groupID},
	}
	res, err := p.listGroupRequestLogs(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list group request logs with group_id: %d failed", groupID)
	}

	return res, nil
}

func (p *postgres) ListGroupRequestLogsBySender(ses storage.Session, sender string) ([]*entity.GroupRequestLog, error) {
	w := &entity.Where{
		FieldNames:  []string{"sender"},
		FieldValues: []any{sender},
	}
	res, err := p.listGroupRequestLogs(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list group request logs with sender: %s failed", sender)
	}

	return res, nil
}

func (p *postgres) GetGroupRequestLog(ses storage.Session, id int64) (*entity.GroupRequestLog, error) {
	w := &entity.Where{
		FieldNames:  []string{"id"},
		FieldValues: []any{id},
	}
	res, err := p.listGroupRequestLogs(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "get group request log with id: %d failed", id)
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no group request log id: %d found", id)
	}

	return res[0], nil
}

func (p *postgres) GetPendingGroupRequestLog(ses storage.Session, groupID int64, sender string) (*entity.GroupRequestLog, error) {
	w := &entity.Where{
		FieldNames:  []string{"group_id", "sender", "status"},
		FieldValues: []any{groupID, sender, entity.LogsStatusPending},
	}
	res, err := p.listGroupRequestLogs(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "get pending group request log with group_id: %d and sender: %s failed", groupID, sender)
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no pending group request log with group_id: %d and sender: %s found", groupID, sender)
	}

	return res[0], nil
}

func (p *postgres) GetPendingGroupRequestLogForUpdate(ses storage.Session, groupID int64, sender string) (*entity.GroupRequestLog, error) {
	sqlstr := rebind(`SELECT * 
                  FROM "group_request_log" 
                  WHERE group_id = ? 
                  AND sender = ?
                  AND status = ?
                  FOR UPDATE;`)

	var res entity.GroupRequestLog
	var err error
	err = ses.QueryRow(sqlstr, groupID, sender, entity.LogsStatusPending).Scan(
		&res.ID, &res.GroupID, &res.Sender, &res.Status, &res.CreatedAt,
	)
	if err != nil {
		return nil, wrapPGErrorf(err, "get pending group request log for update with group_id: %d and sender: %s failed", groupID, sender)
	}

	return &res, nil
}

func (p *postgres) GetGroupRequestLogByIDForUpdate(ses storage.Session, id int64) (*entity.GroupRequestLog, error) {
	sqlstr := rebind(`SELECT * FROM "group_request_log" WHERE id = ? FOR UPDATE;`)

	var res entity.GroupRequestLog
	var err error
	err = ses.QueryRow(sqlstr, id).Scan(
		&res.ID, &res.GroupID, &res.Sender, &res.Status, &res.CreatedAt,
	)
	if err != nil {
		return nil, wrapPGErrorf(err, "get group request log for update with id: %d failed", id)
	}

	return &res, nil
}

func (p *postgres) UpdateGroupRequestLogStatus(ses storage.Session, id int64, status entity.LogsStatus) error {
	sqlstr := rebind(`UPDATE "group_request_log" SET status = ? WHERE id = ?;`)
	_, err := ses.Exec(sqlstr, status, id)
	if err != nil {
		return wrapPGErrorf(err, "update group request log status with id: %d to %s failed", id, status.String())
	}

	return nil
}
