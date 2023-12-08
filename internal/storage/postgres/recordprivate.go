package postgres

import (
	"fmt"
	"sort"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/storage"
)

func (p *postgres) InsertRecordPrivate(ses storage.Session, i *entity.RecordPrivate) (int64, error) {
	sqlstr := rebind(`INSERT INTO "record_private" 
                  (unique_id, content, sender, receiver)
                  VALUES
                  (?, ?, ?, ?)
                  RETURNING id;`)
	args := []any{
		uniqueID(i.Sender, i.Receiver),
		i.Content,
		i.Sender,
		i.Receiver,
	}

	var id int64
	var err error
	err = ses.QueryRow(sqlstr, args...).Scan(&id)
	if err != nil {
		return 0, wrapPGErrorf(err, "failed to insert record_private")
	}

	return id, nil
}

func (p *postgres) listRecordPrivates(ses storage.Session, where *entity.Where) ([]*entity.RecordPrivate, error) {
	projection := []string{
		"id",
		"content",
		"sender",
		"receiver",
		"created_at",
	}

	var args []any
	sqlstr := fmt.Sprintf(`SELECT %s FROM "record_private"`, strings.Join(projection, ", "))
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
		return nil, wrapPGErrorf(err, "failed to list record_private")
	}
	defer rows.Close()

	var res []*entity.RecordPrivate
	for rows.Next() {
		r := entity.RecordPrivate{}
		if err = rows.Scan(&r.ID, &r.Content, &r.Sender, &r.Receiver, &r.CreatedAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan record_private")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) ListRecordPrivatesByParty(ses storage.Session, subject1, subject2 string) ([]*entity.RecordPrivate, error) {
	w := &entity.Where{
		FieldNames:  []string{"unique_id"},
		FieldValues: []any{uniqueID(subject1, subject2)},
	}

	res, err := p.listRecordPrivates(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list record_private with party: %s, %s failed", subject1, subject2)
	}

	return res, nil
}

// 无论subject1和subject2交换与否，保证它们的unique_id一致
func uniqueID(subject1, subject2 string) int64 {
	strs := []string{subject1, subject2}
	sort.Strings(strs)
	uniqueIDStr := strings.Join(strs, "-")
	return hashCode(uniqueIDStr)
}
