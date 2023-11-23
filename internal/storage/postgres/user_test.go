package postgres

import (
	"context"
	"github.com/google/uuid"

	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/storage"
)

func (s *postgresSuite) TestUser() {
	ses, err := s.storage.NewSession(context.Background())
	s.Require().Nil(err)
	ses, err = ses.Begin()
	s.Require().Nil(err)
	defer ses.Rollback()

	subject := uuid.NewString()
	u := &entity.User{
		Subject:  subject,
		Nickname: "foo_nick",
		Username: "foo_name",
		Password: "foo_pw",
		Phone:    "foo_phone",
	}
	err = s.storage.InsertUser(ses, u)
	s.Require().Nil(err)

	got, err := s.storage.GetUserBySubject(ses, subject)
	s.Require().Nil(err)
	s.Require().Equal(u.Nickname, got.Nickname)
	s.Require().Equal(u.Username, got.Username)
	s.Require().Equal(u.Password, got.Password)
	s.Require().Equal(u.Phone, got.Phone)

	err = s.storage.DeleteUser(ses, subject)
	s.Require().Nil(err)

	got, err = s.storage.GetUserBySubject(ses, subject)
	s.Require().Equal(errors.Code(err), errors.NotFound)
}

func (s *postgresSuite) addUser(ses storage.Session) *entity.User {
	u := &entity.User{
		Subject:  uuid.NewString(),
		Nickname: "foo_nick",
		Username: "foo_name",
		Password: "foo_pw",
		Phone:    "foo_phone",
	}
	err := s.storage.InsertUser(ses, u)
	s.Require().Nil(err)

	return u
}
