package postgres

import (
	"fmt"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/storage"
)

func (p *postgres) InsertGroupInvitationLog(ses storage.Session, i *entity.GroupInvitationLog) error {
	sqlstr := rebind(`INSERT INTO "group_invitation_log" 
                  (group_id, sender, receiver, status)
                  VALUES
                  (?, ?, ?, ?);`)
	args := []any{
		i.GroupID,
		i.Sender,
		i.Receiver,
		i.Status,
	}

	var err error
	_, err = ses.Exec(sqlstr, args...)
	if err != nil {
		return wrapPGErrorf(err, "failed to insert group invitation log")
	}

	return nil
}

func (p *postgres) listGroupInvitationLogs(ses storage.Session, where *entity.Where) ([]*entity.GroupInvitationLog, error) {
	projection := []string{
		"id",
		"group_id",
		"sender",
		"receiver",
		"status",
		"created_at",
	}

	var args []any
	sqlstr := fmt.Sprintf(`SELECT %s FROM "group_invitation_log"`, strings.Join(projection, ", "))
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
		return nil, wrapPGErrorf(err, "failed to list group invitation logs")
	}
	defer rows.Close()

	var res []*entity.GroupInvitationLog
	for rows.Next() {
		r := entity.GroupInvitationLog{}
		if err = rows.Scan(&r.ID, &r.GroupID, &r.Sender, &r.Receiver, &r.Status, &r.CreatedAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan group invitation log")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) ListGroupInvitationLogsByReceiver(ses storage.Session, receiver string) ([]*entity.GroupInvitationLog, error) {
	w := &entity.Where{
		FieldNames:  []string{"receiver"},
		FieldValues: []any{receiver},
	}
	res, err := p.listGroupInvitationLogs(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list group invitation logs with receiver: %s failed", receiver)
	}

	return res, nil
}

func (p *postgres) ListGroupInvitationLogsBySender(ses storage.Session, sender string) ([]*entity.GroupInvitationLog, error) {
	w := &entity.Where{
		FieldNames:  []string{"sender"},
		FieldValues: []any{sender},
	}
	res, err := p.listGroupInvitationLogs(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list group invitation logs with sender: %s failed", sender)
	}

	return res, nil
}

func (p *postgres) GetGroupInvitationLog(ses storage.Session, id int64) (*entity.GroupInvitationLog, error) {
	w := &entity.Where{
		FieldNames:  []string{"id"},
		FieldValues: []any{id},
	}
	res, err := p.listGroupInvitationLogs(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "get group invitation log with id: %d failed", id)
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no group invitation log id: %d found", id)
	}

	return res[0], nil
}

func (p *postgres) GetPendingGroupInvitationLog(ses storage.Session, groupID int64, receiver string) (*entity.GroupInvitationLog, error) {
	w := &entity.Where{
		FieldNames:  []string{"group_id", "receiver", "status"},
		FieldValues: []any{groupID, receiver, entity.LogsStatusPending},
	}
	res, err := p.listGroupInvitationLogs(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "get pending group invitation log with group_id: %d and receiver: %s failed", groupID, receiver)
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no pending group invitation log with group_id: %d and receiver: %s found", groupID, receiver)
	}

	return res[0], nil
}

func (p *postgres) GetPendingGroupInvitationLogForUpdate(ses storage.Session, groupID int64, receiver string) (*entity.GroupInvitationLog, error) {
	sqlstr := rebind(`SELECT * 
                  FROM "group_invitation_log" 
                  WHERE group_id = ? 
                  AND receiver = ?
                  AND status = ?
                  FOR UPDATE;`)

	var res entity.GroupInvitationLog
	var err error
	err = ses.QueryRow(sqlstr, groupID, receiver, entity.LogsStatusPending).Scan(
		&res.ID, &res.GroupID, &res.Sender, &res.Receiver, &res.Status, &res.CreatedAt,
	)
	if err != nil {
		return nil, wrapPGErrorf(err, "get pending group invitation log for update with group_id: %d and receiver: %s failed", groupID, receiver)
	}

	return &res, nil
}

func (p *postgres) GetGroupInvitationLogByIDForUpdate(ses storage.Session, id int64) (*entity.GroupInvitationLog, error) {
	sqlstr := rebind(`SELECT * FROM "group_invitation_log" WHERE id = ? FOR UPDATE;`)

	var res entity.GroupInvitationLog
	var err error
	err = ses.QueryRow(sqlstr, id).Scan(
		&res.ID, &res.GroupID, &res.Sender, &res.Receiver, &res.Status, &res.CreatedAt,
	)
	if err != nil {
		return nil, wrapPGErrorf(err, "get group invitation log for update with id: %d failed", id)
	}

	return &res, nil
}

func (p *postgres) UpdateGroupInvitationLogStatus(ses storage.Session, id int64, status entity.LogsStatus) error {
	sqlstr := rebind(`UPDATE "group_invitation_log" SET status = ? WHERE id = ?;`)
	_, err := ses.Exec(sqlstr, status, id)
	if err != nil {
		return wrapPGErrorf(err, "update group invitation log status with id: %d to %s failed", id, status.String())
	}

	return nil
}
