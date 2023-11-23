package main

import (
	"context"
	"log"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"fangaoxs.com/go-chat/server"

	"golang.org/x/sync/errgroup"
)

func main() {
	env, err := environment.Get()
	if err != nil {
		log.Fatalf("init env failed: %v", err)
		return
	}

	logging := logger.New(env)

	s, err := server.New(env, logging)
	if err != nil {
		log.Fatalf("init server failed: %v", err)
		return
	}

	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.Run(ctx)
	})

	if err = g.Wait(); err != nil {
		log.Fatal(err)
		return
	}
}
