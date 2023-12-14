package postgres

import (
	"fmt"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/storage"
)

func (p *postgres) InsertUser(ses storage.Session, i *entity.User) error {
	sqlstr := rebind(`INSERT INTO "user" 
                  (subject, nickname, username, password, phone)
                  VALUES
                  (?, ?, ?, ?, ?);`)
	args := []any{
		i.Subject,
		i.Nickname,
		i.Username,
		i.Password,
		i.Phone,
	}

	var err error
	_, err = ses.Exec(sqlstr, args...)
	if err != nil {
		return wrapPGErrorf(err, "failed to insert user")
	}

	return nil
}

func (p *postgres) listUsers(ses storage.Session, where *entity.Where) ([]*entity.User, error) {
	projection := []string{
		"subject",
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
		if err = rows.Scan(&r.Subject, &r.Nickname, &r.Username, &r.Password, &r.Phone, &r.CreatedAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan user")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) ListAllUsers(ses storage.Session) ([]*entity.User, error) {
	res, err := p.listUsers(ses, nil)
	if err != nil {
		return nil, wrapPGErrorf(err, "list all users failed")
	}

	return res, nil
}

func (p *postgres) GetUserBySubject(ses storage.Session, subject string) (*entity.User, error) {
	w := &entity.Where{
		FieldNames:  []string{"subject"},
		FieldValues: []any{subject},
	}
	users, err := p.listUsers(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "get user with subject: %s failed", subject)
	}
	if len(users) == 0 {
		return nil, errors.Newf(errors.NotFound, nil, "no user with subject: %s found", subject)
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

func (p *postgres) DeleteUser(ses storage.Session, subject string) error {
	sqlstr := rebind(`DELETE FROM "user" WHERE subject = ?;`)
	if _, err := ses.Exec(sqlstr, subject); err != nil {
		return wrapPGErrorf(err, "delete user with subject: %s failed", subject)
	}

	return nil
}

func (p *postgres) InsertFriendship(ses storage.Session, userSubject, friendSubject string) error {
	sqlstr := rebind(`INSERT INTO "friendship" 
                  (user_subject, friend_subject)
                  VALUES
                  (?, ?);`)
	args := []any{
		userSubject,
		friendSubject,
	}

	var err error
	_, err = ses.Exec(sqlstr, args...)
	if err != nil {
		return wrapPGErrorf(err, "failed to insert friendship")
	}

	return nil
}

func (p *postgres) listFriendships(ses storage.Session, where *entity.Where) ([]*entity.Friendship, error) {
	projection := []string{
		"user_subject",
		"friend_subject",
		"created_at",
	}

	var args []any
	sqlstr := fmt.Sprintf(`SELECT %s FROM "friendship"`, strings.Join(projection, ", "))
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
		return nil, wrapPGErrorf(err, "failed to list friendships")
	}
	defer rows.Close()

	var res []*entity.Friendship
	for rows.Next() {
		r := entity.Friendship{}
		if err = rows.Scan(&r.UserSubject, &r.FriendSubject, &r.CreatedAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan friendship")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) IsFriendOfUser(ses storage.Session, userSubject, friendSubject string) (bool, error) {
	w := &entity.Where{
		FieldNames:  []string{"user_subject", "friend_subject"},
		FieldValues: []any{userSubject, friendSubject},
	}
	res, err := p.listFriendships(ses, w)
	if err != nil {
		return false, wrapPGErrorf(err, "is friend: %s of user: %s failed", friendSubject, userSubject)
	}
	if len(res) == 0 {
		return false, nil
	}

	return true, nil
}

func (p *postgres) ListFriendshipsByUserSubject(ses storage.Session, userSubject string) ([]*entity.Friendship, error) {
	w := &entity.Where{
		FieldNames:  []string{"user_subject"},
		FieldValues: []any{userSubject},
	}

	res, err := p.listFriendships(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list friendships with user_subject: %s failed", userSubject)
	}

	return res, nil
}

func (p *postgres) DeleteFriendship(ses storage.Session, userSubject, friendSubject string) error {
	sqlstr := rebind(`DELETE FROM "friendship" WHERE user_subject = ? AND friend_subject = ?;`)
	args := []any{
		userSubject,
		friendSubject,
	}
	if _, err := ses.Exec(sqlstr, args...); err != nil {
		return wrapPGErrorf(err, "delete friendship with user_subject: %s and friend_subject: %s failed", userSubject, friendSubject)
	}

	return nil
}
