package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type app struct {
	srv    *http.Server
	ctx    context.Context
	cancel func()
}

var NilServerError = errors.New("server is nil")

func New(server *http.Server) (*app, error) {
	if server == nil {
		return nil, NilServerError
	}

	ctx, cancel := context.WithCancel(context.TODO())
	return &app{
		srv:    server,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (a *app) Run() error {
	g, ctx := errgroup.WithContext(a.ctx)
	g.Go(func() error {
		<-ctx.Done()
		return a.srv.Shutdown(ctx)
	})
	g.Go(func() error {
		return a.srv.ListenAndServe()
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				a.Stop()
			}
		}
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *app) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}
