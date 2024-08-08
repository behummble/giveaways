package botapp

import (
	"log/slog"

	router "github.com/behummble/giveaway/internal/handlers/http/bot"
	"github.com/behummble/giveaway/internal/service/giveaway"
	"github.com/behummble/giveaway/internal/service/bot"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	log *slog.Logger
	tgbot *telebot.Bot
	giveawayService *giveawayservice.Giveaway
	router *router.Bot
}

func New(log *slog.Logger, botService *botservice.Bot, tgBot *telebot.Bot) (*Bot, error) {
	router, err := router.New(log, botService, tgBot)
	if err != nil {
		return nil, err
	}
	return &Bot{
		log: log,
		tgbot: tgBot,
		botService: botService,
		router: router,
	}, nil
}
