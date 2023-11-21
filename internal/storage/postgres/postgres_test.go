package postgres

import (
	"testing"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/storage"

	"github.com/stretchr/testify/suite"
)

type postgresSuite struct {
	suite.Suite
	storage storage.Storage
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(postgresSuite))
}

// 全部测试开始运行前调用（只执行一次）
func (s *postgresSuite) SetupSuite() {
	env, _ := environment.Get()
	pg, err := New(env)
	s.Require().Nil(err)
	s.storage = pg
}

// 全部测试运行完成后调用（只执行一次）
func (s *postgresSuite) TearDownSuite() {}

// 每个测试开始运行前调用
func (s *postgresSuite) SetupTest() {}

// 每个测试运行完成后调用
func (s *postgresSuite) TearDownTest() {}

func (s *postgresSuite) TestInit() {} // 测试SetupSuite
