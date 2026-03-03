package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"stock/internal/logger"
)

type Runner interface {
	Run() error
	Stop() error
}

func RunApp(app Runner) error {
	done := make(chan error, 1)

	// Use sync.Once to ensure Stop() is only called once and done is sent only once
	var stopOnce sync.Once
	var stopErr error

	stopFunc := func() error {
		stopOnce.Do(func() {
			stopErr = app.Stop()
			// Send to done channel only once
			done <- stopErr
		})
		return stopErr
	}

	cleanupOnInterrupt(stopFunc)
	recoveredRun(done, app, stopFunc)
	return <-done
}

func cleanupOnInterrupt(stopFunc func() error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-c
		logger.Default.Info("got signal", "signal", s.String())
		stopFunc()
	}()
}

func recoveredRun(done chan error, app Runner, stopFunc func() error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Default.Error("panic recovered", "panic", err, "stack", string(debug.Stack()))
			stopFunc()
			return
		}
	}()

	go func() {
		server := &http.Server{
			Addr:         ":6060",
			Handler:      nil,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  2 * time.Minute,
		}

		err := server.ListenAndServe()
		if err != nil {
			logger.Default.Error("failed to listen and serve", "error", err)
			done <- err
		}
	}()

	err := app.Run()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Default.Error("failed to run app", "error", err)
			stopFunc()
			return
		}
	}

	stopFunc()
}
