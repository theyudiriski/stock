package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Server) telegramWebhook(w http.ResponseWriter, r *http.Request) {
	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		s.logger.Error("failed to decode telegram update", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.handleUpdate(update)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		chatID := update.Message.Chat.ID

		s.logger.Info("telegram message received",
			"from", update.Message.From.UserName,
			"text", update.Message.Text,
		)

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Hello, %s!", update.Message.From.UserName))
		s.bot.Send(msg)
	}
}
