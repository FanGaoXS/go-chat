package postgres

import (
	"fangaoxs.com/go-chat/internal/entity"
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
		return 0,
	}

	return 0, nil
}

func (p *postgres) GetUserByID(ses storage.Session, id int64) (*entity.User, error) {
	return nil, nil
}
