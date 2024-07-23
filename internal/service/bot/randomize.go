package botservice

import (
	"log/slog"
	"context"

	"gopkg.in/telebot.v3"
	"github.com/robfig/cron/v3"
	"github.com/behummble/giveaway/internal/entity"
)

type DB interface {
	AddParticipant(ctx context.Context, participant entity.Participant) (int, error)
	Participant(ctx context.Context, id int) (entity.Participant, error)
	AddGiveaway(ctx context.Context, giveaway entity.Giveaway) (int, error)
	Giveaway(ctx context.Context, id int) (entity.Giveaway, error)
	UpdateGiveaway(ctx context.Context, giveaway entity.Giveaway) (entity.Giveaway, error)
	DeleteGiveaway(ctx context.Context, id int) error
	Winners(ctx context.Context, id int) ([]entity.Winner, error)
	Client(ctx context.Context, key string) (entity.Client, error)
	AddClient(ctx context.Context, client entity.Client) (int, error)
	DeleteClient(ctx context.Context, id int) error
}

type Bot struct {
	log *slog.Logger
	tgbot *telebot.Bot
	db DB
	cron *cron.Cron
}

func New(log *slog.Logger, db DB, tgbot *telebot.Bot, cron *cron.Cron) *Bot {
	return &Bot{
		log: log,
		tgbot: tgbot,
		db: db,
		cron: cron,
	}
}

func(bot *Bot) NotifyWinner(ctx context.Context, id int64) error {

}

func(bot *Bot) PublishResults(ctx context.Context, messageID int64, winners []int64) error {
	
}

func(bot *Bot) CancelGiveaway(ctx context.Context, messageID int64) error {
	
}