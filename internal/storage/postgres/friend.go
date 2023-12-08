package postgres

import (
	"fmt"
	"strings"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/storage"
)

func (p *postgres) InsertFriendship(ses storage.Session, i *entity.Friendship) error {
	sqlstr := rebind(`INSERT INTO "friendship" 
                  (user_subject, friend_subject)
                  VALUES
                  (?, ?);`)
	args := []any{
		i.UserSubject,
		i.FriendSubject,
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

func (p *postgres) InsertFriendRequest(ses storage.Session, i *entity.FriendRequest) error {
	sqlstr := rebind(`INSERT INTO "friend_request" 
                  (sender, receiver, status)
                  VALUES
                  (?, ?, ?);`)
	args := []any{
		i.Sender,
		i.Receiver,
		i.Status,
	}

	var err error
	_, err = ses.Exec(sqlstr, args...)
	if err != nil {
		return wrapPGErrorf(err, "failed to insert friend request")
	}

	return nil
}

func (p *postgres) listFriendRequests(ses storage.Session, where *entity.Where) ([]*entity.FriendRequest, error) {
	projection := []string{
		"sender",
		"receiver",
		"status",
		"created_at",
	}

	var args []any
	sqlstr := fmt.Sprintf(`SELECT %s FROM "friend_request"`, strings.Join(projection, ", "))
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
		return nil, wrapPGErrorf(err, "failed to list friend requests")
	}
	defer rows.Close()

	var res []*entity.FriendRequest
	for rows.Next() {
		r := entity.FriendRequest{}
		if err = rows.Scan(&r.Sender, &r.Receiver, &r.Status, &r.CreatedAt); err != nil {
			return nil, wrapPGErrorf(err, "failed to scan friend request")
		}
		res = append(res, &r)
	}

	return res, nil
}

func (p *postgres) ListFriendRequestsBySender(ses storage.Session, sender string) ([]*entity.FriendRequest, error) {
	w := &entity.Where{
		FieldNames:  []string{"sender"},
		FieldValues: []any{sender},
	}

	res, err := p.listFriendRequests(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list friend requests with sender: %s failed", sender)
	}

	return res, nil
}

func (p *postgres) ListFriendRequestsByReceiver(ses storage.Session, receiver string) ([]*entity.FriendRequest, error) {
	w := &entity.Where{
		FieldNames:  []string{"receiver"},
		FieldValues: []any{receiver},
	}

	res, err := p.listFriendRequests(ses, w)
	if err != nil {
		return nil, wrapPGErrorf(err, "list friend requests with receiver: %s failed", receiver)
	}

	return res, nil
}

func (p *postgres) GetFriendRequestForUpdate(ses storage.Session, sender, receiver string) (*entity.FriendRequest, error) {
	sqlstr := rebind(`SELECT * FROM "friend_request WHERE sender = ? AND receiver = ? FOR UPDATE;"`)

	var res *entity.FriendRequest
	var err error
	err = ses.QueryRow(sqlstr, sender, receiver).Scan(
		res.Sender, res.Receiver, res.Status, res.CreatedAt,
	)
	if err != nil {
		return nil, wrapPGErrorf(err, "get friend request for update with sender: %s and receiver: %s failed", sender, receiver)
	}

	return res, nil
}
