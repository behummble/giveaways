package app

import (
	"log/slog"
	"time"
	"github.com/behummble/giveaway/internal/config"
	"github.com/behummble/giveaway/internal/repo/db/postgres"
	"github.com/behummble/giveaway/internal/app/giveaway"
	"github.com/behummble/giveaway/internal/app/bot"
	"github.com/behummble/giveaway/internal/service/bot"
	"gopkg.in/telebot.v3"
	"github.com/robfig/cron/v3"
)

type App struct {
	Bot *botapp.Bot
	Giveaway *giveawayapp.Giveaway
}

func New(log *slog.Logger, config *config.Config) *App {
	
	db, err := newDB(
		log, 
		config.DB.Username, 
		config.DB.Host, 
		config.DB.Username, 
		config.DB.Password,
		config.DB.Port)

	if err != nil {
		panic("Can`t initialize database")
	}

	tgBot, err := newTgBot(config.Bot.Token, config.Bot.UpdateTimeout)
	if err != nil {
		panic("Can`t initialize telegram bot" + err.Error())
	}

	cron, err := newScheduler()
	if err != nil {
		panic("Can`t load timeZone" + err.Error())
	}

	botService := botservice.New(log, db, tgBot, cron)

	botApp, err := botapp.New(log, db, tgBot)
	if err != nil {
		panic("Can`t init bot" + err.Error())
	}

	giveawayApp, err := giveawayapp.New(log, db, botService, cron)

	return &App{
		Bot: botApp,
		Giveaway: giveawayApp,
	}

}

func newDB(log *slog.Logger, username, password, dbname, host string, port int) (postgres.Postgres, error) {
	return postgres.New(
		log,
		username,
		password,
		dbname,
		host,
		port,
	)
}

func newTgBot(token string, timeout int) (*telebot.Bot, error) {
	tgBot, err := telebot.NewBot(
			telebot.Settings{
			Token: token,
			Poller: &telebot.LongPoller{Timeout: time.Second * time.Duration(timeout)},
			},
		)
	return tgBot, err
}

func newScheduler() (*cron.Cron, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, err
	}
	return cron.New(cron.WithLocation(loc)), nil
}