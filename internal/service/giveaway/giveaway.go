package giveawayservice

import "log/slog"

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

type Giveaway struct {
	log *slog.Logger
	db DB
}

func New(log *slog.Logger, db DB) *Giveaway {
	return nil
}

func(g *Giveaway) ClientExist(key string) bool {
	return g.db.Client(key)
}