package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"stock/config"
	"stock/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	server *http.Server
	logger *slog.Logger
	bot    *tgbotapi.BotAPI
}

func New() *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	configServer := config.LoadServer()
	configTelegram := config.LoadTelegram()

	logger.Init()
	log := logger.Default

	bot, err := tgbotapi.NewBotAPI(configTelegram.Bot.Token)
	if err != nil {
		panic(fmt.Sprintf("failed to init telegram bot: %v", err))
	}
	log.Info("telegram bot authorized", "username", bot.Self.UserName)

	s := &Server{
		bot:    bot,
		logger: log,
		server: &http.Server{
			Addr:         configServer.Addr,
			Handler:      r,
			ReadTimeout:  configServer.Timeout.Read,
			WriteTimeout: configServer.Timeout.Write,
			IdleTimeout:  configServer.Timeout.Idle,
		},
	}

	r.Post("/telegram/webhook", s.telegramWebhook)

	return s
}

func (s *Server) Run() error {
	s.logger.Info("http server started", "addr", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("failed to listen and serve", "error", err)
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
