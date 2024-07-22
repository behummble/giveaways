package botservice

import (
	"log/slog"

	"gopkg.in/telebot.v3"
	"github.com/robfig/cron/v3"
)

type DB interface {
	AddParticipant()
	AlredyParticipant()
	AddGiveaway()
	Giveaway()
	UpdateGiveaway()
	DeleteGiveaway()
	Winners()
	Client(key string) bool
}

type Bot struct {
	log *slog.Logger
	tgbot *telebot.Bot
	db DB
	cron *cron.Cron
}

func New(log *slog.Logger, db DB, tgbot *telebot.Bot, cron *cron.Cron) *Bot {
	return &Bot{
		log: log,
		tgbot: tgbot,
		db: db,
		cron: cron,
	}
}
