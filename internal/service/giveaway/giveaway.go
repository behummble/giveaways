package giveawayservice

import (
	"context"
	"log/slog"
	"crypto/rand"

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

type Bot interface {
	NotifyWinner(ctx context.Context, id int64) error
	PublishResults(ctx context.Context, messageID int64, winners []int64) error
	CancelGiveaway(ctx context.Context, messageID int64) error
}

type Giveaway struct {
	log *slog.Logger
	db DB
	bot Bot
}

func New(log *slog.Logger, db DB, bot Bot) *Giveaway {
	return &Giveaway{
		log: log,
		db: db,
		bot: bot,
	}
}

func(giveaway *Giveaway) ClientExist(key string) bool {
	_, err := giveaway.db.Client(context.Background(), key)
	
	return err == nil
}

func(giveaway *Giveaway) AddGiveaway() (int, error) {

}

func(giveaway *Giveaway) Giveaway() (entity.Giveaway, error) {
	
}

func(giveaway *Giveaway) UpdateGiveaway() (entity.Giveaway, error) {
	
}

func(giveaway *Giveaway) DeleteGiveaway() error {
	
}

func(giveaway *Giveaway) AddParticipant() (int, error) {
	
}

func(giveaway *Giveaway) Participant() (entity.Participant, error) {
	
}

func(giveaway *Giveaway) Winners() ([]entity.Winner, error) {
	
}

func(giveaway *Giveaway) Client() (entity.Client, error) {
	
}

func(giveaway *Giveaway) AddClient() (int, error) {
	
}
