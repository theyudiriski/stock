package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
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
			log.Printf("[runner] calling app.Stop()")
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
		log.Printf("[runner] got signal %v", s)
		stopFunc()
	}()
}

func recoveredRun(done chan error, app Runner, stopFunc func() error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
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
			log.Printf("[runner] got error after ListenAndServe %v", err)
			done <- err
		}
	}()

	err := app.Run()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("[runner] got error after Run %v", err)
			stopFunc()
			return
		}
	}

	stopFunc()
}
