package postgres

import (
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fmt"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
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

func (p *postgres) ListGroupByCreatedBy(ses storage.Session, createdBy string) ([]*entity.Group, error) {
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
