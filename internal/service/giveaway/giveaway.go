package giveawayservice

import (
	"log/slog"
	"context"
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

type Bot interface {
	NotifyWinner(ctx context.Context, id int64)
	PublishResults(ctx context.Context, messageID int64, winners []int64)
}

type Giveaway struct {
	log *slog.Logger
	db DB
	bot Bot
}

func New(log *slog.Logger, db DB, bot Bot) *Giveaway {
	return &Giveaway{
		log: log,
		db: db,
		bot: bot,
	}
}

func(g *Giveaway) ClientExist(key string) bool {
	return g.db.Client(key)
}