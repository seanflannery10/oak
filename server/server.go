package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/seanflannery10/ossa/log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	*http.Server
	wg sync.WaitGroup
}

func New(addr string, handler http.Handler) *Server {
	return &Server{
		&http.Server{
			Addr:         addr,
			Handler:      handler,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		sync.WaitGroup{},
	}
}

func (s *Server) Run() error {
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		log.Info("caught signal %s", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		log.Info("completing background tasks on %s", s.Addr)

		s.wg.Wait()
		shutdownError <- nil
	}()

	log.Info("starting server on %s", s.Addr)

	err := s.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Info("server stopped on %s", s.Addr)

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
