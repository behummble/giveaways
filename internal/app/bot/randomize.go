package botapp

import (
	"log/slog"

	router "github.com/behummble/giveaway/internal/handlers/http/bot"
	"github.com/behummble/giveaway/internal/service/bot"
	"gopkg.in/telebot.v3"
)

type BotRouter interface {

}

type Bot struct {
	log *slog.Logger
	tgbot *telebot.Bot
	botService *botservice.Bot
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
