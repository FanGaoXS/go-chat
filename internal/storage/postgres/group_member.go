package postgres

import (
	"fmt"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/storage"
)

func (p *postgres) InsertGroupMember(ses storage.Session, userSubject string, groupID int64) error {
	sqlstr := rebind(`INSERT INTO "group_member"
                  (user_subject, group_id)
                  VALUES
                  (?, ?);`)

	args := []any{
		userSubject,
		groupID,
	}

	_, err := ses.Exec(sqlstr, args...)
	if err != nil {
		return wrapPGErrorf(err, "insert group member with user_subject: %s and group_id: %d failed", userSubject, groupID)
	}

	return nil
}

func (p *postgres) listGroupMembers(ses storage.Session, where *entity.Where) ([]*entity.GroupMember, error) {
	projection := []string{
		"user_subject",
		"group_id",
		"join_at",
	}
	var args []any
	sqlstr := fmt.Sprintf(`SELECT %s FROM "group_member"`, strings.Join(projection, ", "))
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
		return nil, wrapPGErrorf(err, "failed to list group members")
	}
	defer rows.Close()

	var res []*entity.GroupMember
	for rows.Next() {
		r := entity.GroupMember{}
		if err = rows.Scan(&r.UserSubject, &r.GroupID, &r.JoinAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan group member")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) GetGroupMember(ses storage.Session, userSubject string, groupID int64) (*entity.GroupMember, error) {
	w := &entity.Where{
		FieldNames:  []string{"user_subject", "group_id"},
		FieldValues: []any{userSubject, groupID},
	}

	res, err := p.listGroupMembers(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "get group member with user_subject: %s and group_id: %d failed", userSubject, groupID)
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no group member with user_subject: %s and group_id: %d found", userSubject, groupID)
	}

	return res[0], nil
}

func (p *postgres) ListGroupMemberByGroupID(ses storage.Session, groupID int64) ([]*entity.GroupMember, error) {
	w := &entity.Where{
		FieldNames:  []string{"group_id"},
		FieldValues: []any{groupID},
	}

	res, err := p.listGroupMembers(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list group members with group_id: %d failed", groupID)
	}

	return res, nil
}

func (p *postgres) ListGroupMemberByUserSubject(ses storage.Session, userSubject string) ([]*entity.GroupMember, error) {
	w := &entity.Where{
		FieldNames:  []string{"user_subject"},
		FieldValues: []any{userSubject},
	}

	res, err := p.listGroupMembers(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list group members with user_subject: %s failed", userSubject)
	}

	return res, nil
}
