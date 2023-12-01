package postgres

import (
	"fmt"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/storage"
)

func (p *postgres) InsertRecord(ses storage.Session, i *entity.Record) error {
	sqlstr := rebind(`INSERT INTO "record" 
                  (seq_id, "type", metadata, sender, receiver, content)
                  VALUES
                  (?, ?, ?, ?, ?, ?)`)
	args := []any{
		i.SeqID,
		i.Type,
		i.Metadata,
		i.Sender,
		i.Receiver,
		i.Content,
	}

	_, err := ses.Exec(sqlstr, args...)
	if err != nil {
		return wrapPGErrorf(err, "failed to insert record")
	}

	return nil
}

func (p *postgres) listRecords(ses storage.Session, where *entity.Where) ([]*entity.Record, error) {
	projection := []string{
		"seq_id",
		"type",
		"metadata",
		"sender",
		"receiver",
		"content",
		"created_at",
	}

	var args []any
	sqlstr := fmt.Sprintf(`SELECT %s FROM "record"`, strings.Join(projection, ", "))
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
		return nil, wrapPGErrorf(err, "failed to list record")
	}
	defer rows.Close()

	var res []*entity.Record
	for rows.Next() {
		r := entity.Record{}
		if err = rows.Scan(&r.SeqID, &r.Type, &r.Metadata, &r.Sender, &r.Receiver, &r.Content, &r.CreatedAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan record")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) ListBroadcastRecords(ses storage.Session) ([]*entity.Record, error) {
	w := &entity.Where{
		FieldNames:  []string{"type"},
		FieldValues: []any{entity.RecordTypeBroadcast},
	}

	res, err := p.listRecords(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list broadcast records failed")
	}

	return res, nil
}

func (p *postgres) ListGroupRecords(ses storage.Session, metadata string) ([]*entity.Record, error) {
	w := &entity.Where{
		FieldNames:  []string{"type", "metadata"},
		FieldValues: []any{entity.RecordTypeGroup, metadata},
	}

	res, err := p.listRecords(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list group records with metadata: %s failed", metadata)
	}

	return res, nil
}

func (p *postgres) ListPrivateRecords(ses storage.Session, sender, receiver string) ([]*entity.Record, error) {
	w := &entity.Where{
		FieldNames:  []string{"type", "sender", "receiver"},
		FieldValues: []any{entity.RecordTypePrivate, sender, receiver},
	}

	res, err := p.listRecords(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list private records with sender: %s and receiver: %s failed", sender, receiver)
	}

	return res, nil
}
