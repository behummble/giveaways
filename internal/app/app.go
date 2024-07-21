package app

import (
	"log/slog"
	"github.com/behummble/giveaway/internal/config"
)

type App struct {
	Bot int
	Giveaway int
}

func New(log *slog.Logger, config *config.Config) *App {
	return nil
}