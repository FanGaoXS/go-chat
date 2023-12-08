package postgres

import (
	"fmt"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/storage"
)

func (p *postgres) InsertRecordGroup(ses storage.Session, i *entity.RecordGroup) (int64, error) {
	sqlstr := rebind(`INSERT INTO "record_group" 
                  (group_id, content, sender)
                  VALUES
                  (?, ?, ?)
                  RETURNING id;`)
	args := []any{
		i.GroupID,
		i.Content,
		i.Sender,
	}

	var id int64
	var err error
	err = ses.QueryRow(sqlstr, args...).Scan(&id)
	if err != nil {
		return 0, wrapPGErrorf(err, "failed to insert record_group")
	}

	return id, nil
}

func (p *postgres) listRecordGroups(ses storage.Session, where *entity.Where) ([]*entity.RecordGroup, error) {
	projection := []string{
		"id",
		"group_id",
		"content",
		"sender",
		"created_at",
	}

	var args []any
	sqlstr := fmt.Sprintf(`SELECT %s FROM "record_group"`, strings.Join(projection, ", "))
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
		return nil, wrapPGErrorf(err, "failed to list record_group")
	}
	defer rows.Close()

	var res []*entity.RecordGroup
	for rows.Next() {
		r := entity.RecordGroup{}
		if err = rows.Scan(&r.ID, &r.GroupID, &r.Content, &r.Sender, &r.CreatedAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan record_group")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) ListRecordGroupsByGroup(ses storage.Session, groupID int64) ([]*entity.RecordGroup, error) {
	w := &entity.Where{
		FieldNames:  []string{"group_id"},
		FieldValues: []any{groupID},
	}

	res, err := p.listRecordGroups(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list record groups with group: %d failed", groupID)
	}

	return res, nil
}
