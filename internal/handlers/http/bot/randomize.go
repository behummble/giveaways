package bot

import (
	"log/slog"
	"time"
	"gopkg.in/telebot.v3"
)

type BotService interface {

}

type Bot struct {
	log *slog.Logger
	tgbot *telebot.Bot
	botService BotService
}

func New(log *slog.Logger, botServie BotService, token string, timeout int) (*Bot, error) {
	tgBot, err := initialize(token, timeout)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		log: log,
		tgbot: tgBot,
		botService: botServie,
	}

	bot.register()

	return bot, nil
}

func initialize(token string, timeout int) (*telebot.Bot, error) {
	return telebot.NewBot(
		telebot.Settings{
			Token: token,
			Poller: &telebot.LongPoller{Timeout: time.Second * time.Duration(timeout)},
		},
	)
}

func(bot *Bot) register() {
	bot.tgbot.Handle("/list", bot.giveaways)
}

func(bot *Bot) giveaways(upd telebot.Context) error {
	return nil
}