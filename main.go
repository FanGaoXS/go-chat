package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	closec := make(chan os.Signal, 1)
	signal.Notify(closec, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.Run(ctx)
	})

	go func() {
		<-closec
		cancel()
	}()

	if err = g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Println(err)
	}

	log.Println(env.AppName, env.AppVersion, "stopped")
}
