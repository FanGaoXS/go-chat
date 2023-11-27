package postgres

import (
	"fmt"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/storage"
)

func (p *postgres) InsertGroup(ses storage.Session, i *entity.Group) (int64, error) {
	sqlstr := rebind(`INSERT INTO "group" 
                  (name, "type", is_public, created_by)
                  VALUES
                  (?, ?, ?, ?)
                  RETURNING id;`)
	args := []any{
		i.Name,
		i.Type,
		i.IsPublic,
		i.CreatedBy,
	}

	var id int64
	var err error
	err = ses.QueryRow(sqlstr, args...).Scan(&id)
	if err != nil {
		return 0, wrapPGErrorf(err, "failed to insert group")
	}

	return id, nil
}

func (p *postgres) listGroups(ses storage.Session, where *entity.Where) ([]*entity.Group, error) {
	projection := []string{
		"id",
		"name",
		"type",
		"is_public",
		"created_by",
		"created_at",
	}

	var args []any
	sqlstr := fmt.Sprintf(`SELECT %s FROM "group"`, strings.Join(projection, ", "))
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
		return nil, wrapPGErrorf(err, "failed to list groups")
	}
	defer rows.Close()

	var res []*entity.Group
	for rows.Next() {
		r := entity.Group{}
		if err = rows.Scan(&r.ID, &r.Name, &r.Type, &r.IsPublic, &r.CreatedBy, &r.CreatedAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan group")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) GetGroupByID(ses storage.Session, id int64) (*entity.Group, error) {
	w := &entity.Where{
		FieldNames:  []string{"id"},
		FieldValues: []any{id},
	}

	groups, err := p.listGroups(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "get group with id: %d failed", id)
	}
	if len(groups) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no group with id: %d found", id)
	}

	return groups[0], nil
}

func (p *postgres) ListGroupsByCreatedBy(ses storage.Session, createdBy string) ([]*entity.Group, error) {
	w := &entity.Where{
		FieldNames:  []string{"created_by"},
		FieldValues: []any{createdBy},
	}

	groups, err := p.listGroups(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list groups with created_by: %s failed", createdBy)
	}

	return groups, nil
}

func (p *postgres) DeleteGroup(ses storage.Session, id int64) error {
	sqlstr := rebind(`DELETE FROM "group" WHERE id = ?;`)
	if _, err := ses.Exec(sqlstr, id); err != nil {
		return wrapPGErrorf(err, "delete user with subject: %d failed", id)
	}

	return nil
}

func (p *postgres) UpdateGroupIsPublic(ses storage.Session, id int64, isPublic bool) error {
	sqlstr := rebind(`UPDATE "group" 
                  SET is_public = ? 
                  WHERE id = ?;`)

	args := []any{isPublic, id}

	_, err := ses.Exec(sqlstr, args...)
	if err != nil {
		return wrapPGErrorf(err, "update is_public of group with id: %d to %t failed", id, isPublic)
	}

	return nil
}

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

func (p *postgres) ListGroupMembersByGroupID(ses storage.Session, groupID int64) ([]*entity.GroupMember, error) {
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

func (p *postgres) ListGroupMembersByUserSubject(ses storage.Session, userSubject string) ([]*entity.GroupMember, error) {
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
