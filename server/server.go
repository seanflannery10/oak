package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/seanflannery10/oak/log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	wg sync.WaitGroup
}

func New() *Server {
	return &Server{}
}

func (s *Server) Run(addr string, h http.Handler) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      h,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		log.Info("caught signal %s", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		log.Info("completing background tasks on %s", srv.Addr)

		s.wg.Wait()
		shutdownError <- nil
	}()

	log.Info("starting server on %s", srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Info("server stopped on %s", srv.Addr)

	return nil
}

func (s *Server) background(fn func()) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				log.Error(fmt.Errorf("%s", err), nil)
			}
		}()

		fn()
	}()
}
