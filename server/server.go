package server

import (
	"context"
	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"fangaoxs.com/go-chat/server/rest"
	"fangaoxs.com/go-chat/server/websocket"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func New(env environment.Env, logger logger.Logger) (*Server, error) {
	httpServer := gin.New()
	gin.ForceConsoleColor()
	httpServer.Use(gin.Logger())

	server, err := initServer(env, logger, httpServer)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func newServer(
	env environment.Env,
	logger logger.Logger,
	httpServer *gin.Engine,
	authorizer auth.Authorizer,
	user user.User,
	group group.Group,
) (*Server, error) {
	restServer, err := rest.New(env, logger, httpServer, authorizer, user, group)
	if err != nil {
		return nil, err
	}

	wsServer, err := websocket.New(env, logger, httpServer, authorizer, user, group)
	if err != nil {
		return nil, err
	}

	return &Server{
		env:        env,
		logger:     logger,
		restServer: restServer,
		wsServer:   wsServer,
	}, nil
}

type Server struct {
	env    environment.Env
	logger logger.Logger

	restServer *rest.Server
	wsServer   *websocket.Server
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

	g.Go(func() error {
		s.logger.Infof("websocket server listen on %s", s.env.WebsocketListenAddr)
		err := s.wsServer.ListenAndServe()
		if err != nil {
			return err
		}
		s.logger.Infof("websocket server stopped")
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
	s.wsServer.Close()
	return nil
}
