package giveawayapp

import (
	"log/slog"

	"github.com/behummble/giveaway/internal/service/giveaway"
)

type Giveaway struct {
	log *slog.Logger
	giveawayService *giveawayservice.Giveaway
}

func New(log *slog.Logger) *Giveaway {
	return nil
}