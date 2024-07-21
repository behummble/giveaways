package botapp

import (
	"log/slog"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	log *slog.Logger
	tgbot *telebot.Bot
}

func New() *Bot {
	return nil
}