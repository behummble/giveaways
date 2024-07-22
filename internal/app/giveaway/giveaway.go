package giveawayapp

import (
	"log/slog"

	router "github.com/behummble/giveaway/internal/handlers/http/giveaway"
	"github.com/behummble/giveaway/internal/service/giveaway"
	"github.com/robfig/cron/v3"
)

type Giveaway struct {
	log *slog.Logger
	giveawayService *giveawayservice.Giveaway
	router *router.Router
}

func New(log *slog.Logger, db giveawayservice.DB, bot giveawayservice.Bot, cron *cron.Cron) (*Giveaway, error) {
	giveawayService := giveawayservice.New(log, db, bot)
	
	router := router.New(log, giveawayService)
	return &Giveaway{
		log: log,
		giveawayService: giveawayService,
		router: router,
	}, nil
}