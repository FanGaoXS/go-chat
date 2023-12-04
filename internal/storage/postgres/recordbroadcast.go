package postgres

import (
	"fmt"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/storage"
)

func (p *postgres) InsertRecordBroadcast(ses storage.Session, i *entity.RecordBroadcast) (int64, error) {
	sqlstr := rebind(`INSERT INTO "record_broadcast" 
                  (content, sender)
                  VALUES
                  (?, ?)
                  RETURNING id;`)
	args := []any{
		i.Content,
		i.Sender,
	}

	var id int64
	var err error
	err = ses.QueryRow(sqlstr, args...).Scan(&id)
	if err != nil {
		return 0, wrapPGErrorf(err, "failed to insert record_broadcast")
	}

	return id, nil
}

func (p *postgres) listRecordBroadcasts(ses storage.Session, where *entity.Where) ([]*entity.RecordBroadcast, error) {
	projection := []string{
		"id",
		"content",
		"sender",
		"created_at",
	}

	var args []any
	sqlstr := fmt.Sprintf(`SELECT %s FROM "record_broadcast"`, strings.Join(projection, ", "))
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
		return nil, wrapPGErrorf(err, "failed to list record_broadcast")
	}
	defer rows.Close()

	var res []*entity.RecordBroadcast
	for rows.Next() {
		r := entity.RecordBroadcast{}
		if err = rows.Scan(&r.ID, &r.Content, &r.Sender, &r.CreatedAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan record_broadcast")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) ListAllRecordBroadcasts(ses storage.Session) ([]*entity.RecordBroadcast, error) {
	res, err := p.listRecordBroadcasts(ses, nil)
	if err != nil {
		return nil, wrapPGErrorf(err, "list all record broadcast failed")
	}

	return res, nil
}

func (p *postgres) ListRecordBroadcastsBySender(ses storage.Session, sender string) ([]*entity.RecordBroadcast, error) {
	w := &entity.Where{
		FieldNames:  []string{"sender"},
		FieldValues: []any{sender},
	}

	res, err := p.listRecordBroadcasts(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list record_broadcast with sender: %s failed", sender)
	}

	return res, nil
}
