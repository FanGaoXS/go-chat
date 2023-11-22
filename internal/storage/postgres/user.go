package postgres

import (
	"fmt"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/storage"
)

func (p *postgres) InsertUser(ses storage.Session, i *entity.User) (int64, error) {
	sqlstr := rebind(`INSERT INTO "user" 
                  (nickname, username, password, phone)
                  VALUES
                  (?, ?, ?, ?)
                  RETURNING id;`)
	args := []any{
		i.Nickname,
		i.Username,
		i.Password,
		i.Phone,
	}

	var id int64
	var err error
	err = ses.QueryRow(sqlstr, args...).Scan(&id)
	if err != nil {
		return 0, wrapPGErrorf(err, "failed to insert user")
	}

	return id, nil
}

func (p *postgres) listUsers(ses storage.Session, where *entity.Where) ([]*entity.User, error) {
	projection := []string{
		"id",
		"nickname",
		"username",
		"password",
		"phone",
		"created_at",
	}

	var args []any
	sqlstr := fmt.Sprintf(`SELECT %s FROM "user"`, strings.Join(projection, ", "))
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
		return nil, wrapPGErrorf(err, "failed to list users")
	}
	defer rows.Close()

	var res []*entity.User
	for rows.Next() {
		r := entity.User{}
		if err = rows.Scan(&r.ID, &r.Nickname, &r.Username, &r.Password, &r.Phone, &r.CreatedAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan user")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) GetUserByID(ses storage.Session, id int64) (*entity.User, error) {
	w := &entity.Where{
		FieldNames:  []string{"id"},
		FieldValues: []any{id},
	}
	users, err := p.listUsers(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "get user with id: %d failed", id)
	}
	if len(users) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no user with id: %d found", id)
	}

	return users[0], nil
}

func (p *postgres) GetUserBySecret(ses storage.Session, username, password string) (*entity.User, error) {
	w := &entity.Where{
		FieldNames:  []string{"username", "password"},
		FieldValues: []any{username, password},
	}
	res, err := p.listUsers(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "get user with secret failed")
	}
	if len(res) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no user found with username: %s and password: %s", username, password)
	}

	return res[0], nil
}
