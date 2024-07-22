package bot

import (
	"log/slog"
	"gopkg.in/telebot.v3"
)

type BotService interface {

}

type Bot struct {
	log *slog.Logger
	tgbot *telebot.Bot
	botService BotService
}

func New(log *slog.Logger, botServie BotService, tgbot *telebot.Bot) (*Bot, error) {
	bot := &Bot{
		log: log,
		tgbot: tgbot,
		botService: botServie,
	}

	return bot, nil
}

func(bot *Bot) Register() {
	bot.tgbot.Handle("/list", bot.giveaways)
}

func(bot *Bot) giveaways(upd telebot.Context) error {
	return nil
}