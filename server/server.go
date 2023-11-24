package server

import (
	"context"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/groupmember"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"fangaoxs.com/go-chat/server/rest"

	"golang.org/x/sync/errgroup"
)

func New(env environment.Env, logger logger.Logger) (*Server, error) {
	server, err := initServer(env, logger)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func newServer(env environment.Env, logger logger.Logger, user user.User, group group.Group, groupMember groupmember.GroupMember) (*Server, error) {
	restServer, err := rest.New(env, logger, user, group, groupMember)
	if err != nil {
		return nil, err
	}

	return &Server{
		env:        env,
		logger:     logger,
		restServer: restServer,
	}, nil
}

type Server struct {
	env    environment.Env
	logger logger.Logger

	restServer *rest.Server
}

func (s *Server) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		s.logger.Infof("rest server listen on %s", s.env.RestListenAddr)
		err := s.restServer.ListenAndServe()
		if err != nil {
			return err
		}
		s.logger.Infof("rest server stopped")
		return nil
	})

	defer s.Close()

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Close() error {
	s.restServer.Close()
	return nil
}
