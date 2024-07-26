package giveawayservice

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/behummble/giveaway/internal/entity"
	"github.com/robfig/cron/v3"
)

type DB interface {
	AddParticipant(ctx context.Context, participant entity.Participant) (int, error)
	Participant(ctx context.Context, id int) (entity.Participant, error)
	AllParticipants(ctx context.Context, giveawayID int) ([]entity.Participant, error)
	AddGiveaway(ctx context.Context, giveaway entity.Giveaway) (int, error)
	Giveaway(ctx context.Context, id int) (entity.Giveaway, error)
	UpdateGiveaway(ctx context.Context, giveaway entity.Giveaway) (entity.Giveaway, error)
	DeleteGiveaway(ctx context.Context, id int) error
	AddWinners(ctx context.Context, winners []entity.Winner) error
	Winners(ctx context.Context, id int) ([]entity.Winner, error)
	Client(ctx context.Context, key string) (entity.Client, error)
	AddClient(ctx context.Context, client entity.Client) (int, error)
	DeleteClient(ctx context.Context, id int) error
}

type Bot interface {
	ParticipantMetTerms(ctx context.Context, terms []string, userID int64) (bool, error)
	NotifyWinner(ctx context.Context, id int64) error
	SchedulePublication(ctx context.Context, giveawayID int) error
	PublishResults(ctx context.Context, messageID int64, winners []int64) error
	CancelGiveaway(ctx context.Context, messageID int64) error
}

type Giveaway struct {
	log *slog.Logger
	db DB
	bot Bot
	cron *cron.Cron
}

func New(log *slog.Logger, db DB, bot Bot, cron *cron.Cron) *Giveaway {
	return &Giveaway{
		log: log,
		db: db,
		bot: bot,
		cron: cron,
	}
}

func (giveaway *Giveaway) ClientExist(key string) bool {
	_, err := giveaway.db.Client(context.Background(), key)
	
	return err == nil
}

func (giveaway *Giveaway) AddGiveaway(body []byte) (int, error) {
	var giveawayModel entity.Giveaway
	err := json.Unmarshal(body, &giveawayModel)
	if err != nil {
		return -1, err
	}

	id, err := giveaway.db.AddGiveaway(context.Background(), giveawayModel)
	if err != nil {
		return -1, err
	}

	err = giveaway.scheduleGiveaway(id, giveawayModel.EventDate, giveawayModel.PublishDate)

	if err != nil {
		//rollback
		return -1, err
	}
}

func (giveaway *Giveaway) Giveaway(giveawayID int) (entity.Giveaway, error) {
	
}

func (giveaway *Giveaway) UpdateGiveaway(giveawayID int, body []byte) (entity.Giveaway, error) {
	
}

func (giveaway *Giveaway) DeleteGiveaway(giveawayID int) error {
	
}

func (giveaway *Giveaway) AddParticipant(body []byte) (int, error) {
	
}

func (giveaway *Giveaway) Participant(participantID int) (entity.Participant, error) {
	
}

func (giveaway *Giveaway) Winners(giveawayID int) ([]entity.Winner, error) {
	
}

func (giveaway *Giveaway) Client(clientID int) (entity.Client, error) {
	
}

func (giveaway *Giveaway) AddClient(body []byte) (int, error) {
	
}

func (giveaway *Giveaway) scheduleGiveaway(giveawayID int, giveawayDate time.Time, publishDate time.Time) error {
	specGiveaway := fmt.Sprintf(
		"%d %d %d %d %d ? %d",
		giveawayDate.Second(),
		giveawayDate.Minute(),
		giveawayDate.Hour(),
		giveawayDate.Day(),
		giveawayDate.Month(),
		giveawayDate.Year())

	
	idFunc, err := giveaway.cron.AddFunc(
		specGiveaway, startGiveawayFunc(giveaway, giveawayID))
}

func startGiveawayFunc(giveaway *Giveaway, giveawayID int) func() {
	return func () {
		giveaway.startGiveaway(giveawayID)
	}
}

func (giveaway *Giveaway) startGiveaway(giveawayID int) {
	participants, err := giveaway.db.AllParticipants(context.Background(), giveawayID)
	if err != nil {
		giveaway.log.Error(
			fmt.Sprintf("Can`t get participants from db by giveawayID: %d", giveawayID),
			"Error", err)
		return
	}

	giveawayData, err := giveaway.db.Giveaway(context.Background(), giveawayID)
	if err != nil {
		giveaway.log.Error(
			fmt.Sprintf("Can`t get giveaway from db by giveawayID: %d", giveawayID),
			"Error", err)
		return
	}

	winners, err := giveaway.findWinners(participants, giveawayData.Terms, giveawayData.WinnersNumber, giveawayID)
	if err != nil {
		giveaway.log.Error(
			fmt.Sprintf("Can`t pick winners for giveaway: %d", giveawayID),
			"Error", err)
			return
	}

	err = giveaway.db.AddWinners(context.Background(), winners)

}

func (giveaway *Giveaway) findWinners(participants []entity.Participant, terms []string, winnersNumber, giveawayID int) ([]int, error) {
	winners := make([]int, 0, winnersNumber) // переделать возврат на структуру Winner и после каждой итерации сокращать кол-во возможных виннеров

	iterations := 0
	for {
		if len(winners) > winnersNumber || len(participants) == 0 {
			break
		}

		if iterations > 100 {
			return []int{}, fmt.Errorf("Can`t generate random value while find winners in giveaway %d", giveawayID)
		}

		max := big.NewInt(int64(len(participants)))
		winnerID, err := rand.Int(rand.Reader, max)

		if err != nil {
			giveaway.log.Error("Can`t generate random value while find winners", "Error", err)
			iterations++
			continue
		}

		index := int(winnerID.Int64())

		winnerValid, err := giveaway.bot.ParticipantMetTerms(
			context.Background(), 
			terms, 
			participants[index].UserID) 

		if err != nil {
			giveaway.log.Error(
				fmt.Sprintf("Can`t verify winner: %d in giveaway: %d", participants[index].UserID, giveawayID),
				"Error", err)
			iterations++
			continue
		}

		if winnerValid {
			winners = append(winners, index)
		}

		if index == len(participants) - 1 {
			participants = participants[:index]
		} else {
			// может переделать на перенос к последнему элементу и обрезать
			//participants = append(participants[:index], participants[index+1:]...)
			participants[index] = participants[len(participants) - 1]
			participants = participants[:len(participants) - 1]
		}
	}

	return winners, nil
}