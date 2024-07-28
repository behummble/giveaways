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
	AddParticipant(ctx context.Context, participant entity.Participant) (entity.Participant, error)
	Participant(ctx context.Context, id int) (entity.Participant, error)
	AllParticipants(ctx context.Context, giveawayID int) ([]entity.Participant, error)
	DeleteAllGiveawayParticipants(ctx context.Context, giveawayID int) error
	AddGiveaway(ctx context.Context, giveaway entity.Giveaway) (entity.Giveaway, error)
	Giveaway(ctx context.Context, id int) (entity.Giveaway, error)
	UpdateGiveaway(ctx context.Context, giveaway entity.Giveaway) (entity.Giveaway, error)
	DeleteGiveaway(ctx context.Context, id int) error
	AddWinners(ctx context.Context, winners []entity.Winner) error
	Winners(ctx context.Context, giveawayID int) ([]entity.Winner, error)
	Client(ctx context.Context, key string) (entity.Client, error)
	AddClient(ctx context.Context, client entity.Client) (int, error)
	DeleteClient(ctx context.Context, id int) error
}

type Bot interface {
	ParticipantMetTerms(ctx context.Context, terms []string, userID int64) (bool, error)
	NotifyWinners(ctx context.Context, giveawayID int)
	SchedulePublication(ctx context.Context, publishDate time.Time, chatID int64, giveawayID int) (int, error)
	PublishResults(ctx context.Context, giveawayID int)
	CancelGiveaway(ctx context.Context, giveawayID int)
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

func (giveaway *Giveaway) AddGiveaway(body []byte) ([]byte, error) {
	var giveawayModel entity.Giveaway
	err := json.Unmarshal(body, &giveawayModel)
	if err != nil {
		giveaway.log.Error("Can`t parse giveaway data from request body", "Error", err)
		return []byte{}, err
	}

	giveawayModel, err = giveaway.db.AddGiveaway(context.Background(), giveawayModel)
	if err != nil {
		giveaway.log.Error("Can`t add giveaway to db", "Error", err)
		return []byte{}, err
	}

	// планирование розыгрыша и планирование публикации можно распараллелить
	idFunc, err := giveaway.scheduleGiveaway(giveawayModel.ID, giveawayModel.EventDate)

	if err != nil {
		//rollback
		giveaway.log.Error("Can`t schedule giveaway", "Error", err)
		return []byte{}, err
	}

	giveawayModel.SchedulerGiveawayID = idFunc

	idFunc, err = giveaway.bot.SchedulePublication(
		context.Background(), 
		giveawayModel.PublishDate, 
		giveawayModel.GiveawayChatID, 
		giveawayModel.ID)

	if err != nil {
		//rollback
		giveaway.log.Error("Can`t schedule publication giveaway", "Error", err)
		return []byte{}, err
	}

	giveawayModel.SchedulerPublishID = idFunc

	giveawayModel, err = giveaway.db.UpdateGiveaway(context.Background(), giveawayModel)

	if err != nil {
		giveaway.log.Error(fmt.Sprintf("Can`t update giveaway %d in db", giveawayModel.ID), "Error", err)
		return []byte{}, err
	}

	response, err := json.Marshal(giveawayModel)
	if err != nil {
		giveaway.log.Error(
			fmt.Sprintf("Can`t serealize giveaway %d data from db to json", giveawayModel.ID), 
			"Error", err)
	}
	
	return response, err
}

func (giveaway *Giveaway) Giveaway(giveawayID int) ([]byte, error) {
	giveawayModel, err := giveaway.db.Giveaway(context.Background(), giveawayID)
	if err != nil {
		giveaway.log.Error(fmt.Sprintf("Can`t get giveaway %d from db", giveawayID), "Error", err)
		return []byte{}, err
	}

	response, err := json.Marshal(giveawayModel)

	if err != nil {
		giveaway.log.Error(fmt.Sprintf("Can`t serealize giveaway %d data from db to json", giveawayID), "Error", err)
	}
	
	return response, err
}

func (giveaway *Giveaway) UpdateGiveaway(giveawayID int, body []byte) ([]byte, error) {
	var giveawayModel entity.Giveaway
	err := json.Unmarshal(body, &giveawayModel)
	if err != nil {
		giveaway.log.Error(fmt.Sprintf("Can`t parse giveaway %d data from request body", giveawayID), "Error", err)
		return []byte{}, err
	}

	giveawayModel, err = giveaway.db.UpdateGiveaway(context.Background(), giveawayModel)
	if err != nil {
		giveaway.log.Error(fmt.Sprintf("Can`t update giveaway %d in db", giveawayID), "Error", err)
		return []byte{}, err
	}

	response, err := json.Marshal(giveawayModel)
	if err != nil {
		giveaway.log.Error(fmt.Sprintf("Can`t serealize giveaway %d data from db to json", giveawayID), "Error", err)
	}
	
	return response, err
}

func (giveaway *Giveaway) DeleteGiveaway(giveawayID int) error {
	go giveaway.bot.CancelGiveaway(context.Background(), giveawayID)
	err := giveaway.db.DeleteAllGiveawayParticipants(context.Background(), giveawayID)
	if err != nil {
		giveaway.log.Error(
			fmt.Sprintf("Couldn`t delete participants by giveawayID %d from db", giveawayID),
			"Error", err)
	}

	err = giveaway.db.DeleteGiveaway(context.Background(), giveawayID)
	if err != nil {
		giveaway.log.Error(
			fmt.Sprintf("Couldn`t delete giveaway %d from db", giveawayID),
			"Error", err)
	}

	return err
}

func (giveaway *Giveaway) AddParticipant(body []byte) ([]byte, error) {
	var participantModel entity.Participant
	err := json.Unmarshal(body, &participantModel)
	if err != nil {
		giveaway.log.Error("Can`t parse participant data from request body", "Error", err)
		return []byte{}, err
	}

	participantModel, err = giveaway.db.AddParticipant(context.Background(), participantModel)
	if err != nil {
		giveaway.log.Error("Can`t add participant to db", "Error", err)
		return []byte{}, err
	}

	response, err := json.Marshal(participantModel)
	if err != nil {
		giveaway.log.Error(fmt.Sprintf("Can`t serealize participant %d data from db to json", participantModel.ID), "Error", err)
	}
	
	return response, err
}

func (giveaway *Giveaway) Participant(participantID int) ([]byte, error) {
	participantModel, err := giveaway.db.Participant(context.Background(), participantID)
	if err != nil {
		giveaway.log.Error(fmt.Sprintf("Can`t get participant %d from db", participantID), "Error", err)
		return []byte{}, err
	}

	response, err := json.Marshal(participantModel)
	if err != nil {
		giveaway.log.Error(fmt.Sprintf("Can`t serealize participant %d data from db to json", participantID), "Error", err)
	}
	
	return response, err
}

func (giveaway *Giveaway) Winners(giveawayID int) ([]byte, error) {
	winners, err := giveaway.db.Winners(context.Background(), giveawayID)
	if err != nil {
		giveaway.log.Error(fmt.Sprintf("Can`t get winners giveaway %d from db", giveawayID), "Error", err)
		return []byte{}, err
	}

	response, err := json.Marshal(winners)
	if err != nil {
		giveaway.log.Error(fmt.Sprintf("Can`t serealize winners giveaway %d from db to json", giveawayID), "Error", err)
	}
	
	return response, err
}

func (giveaway *Giveaway) Client(clientID int) ([]byte, error) {
	/*clientModel, err := giveaway.db.Client(context.Background(), clientID)
	if err != nil {
		return []byte{}, err
	}

	response, err := json.Marshal(clientModel)
	
	return response, err */
	return []byte{}, nil
}

func (giveaway *Giveaway) AddClient(body []byte) (int, error) {
	/*var clientModel entity.Client
	err := json.Unmarshal(body, &clientModel)
	if err != nil {
		return -1, err
	}

	clientID, err := giveaway.db.AddClient(context.Background(), clientModel)
	if err != nil {
		return -1, err
	}

	response, err := json.Marshal
	
	return response, err */
	return -1, nil
}

func (giveaway *Giveaway) scheduleGiveaway(giveawayID int, giveawayDate time.Time) (int, error) {
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

	return int(idFunc), err
}

func startGiveawayFunc(giveaway *Giveaway, giveawayID int) func() {
	return func () {
		giveaway.startGiveaway(giveawayID)
		go giveaway.bot.PublishResults(context.Background(), giveawayID)
		giveaway.bot.NotifyWinners(context.Background(), giveawayID)
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

func (giveaway *Giveaway) findWinners(participants []entity.Participant, terms []string, winnersNumber, giveawayID int) ([]entity.Winner, error) {
	winners := make([]entity.Winner, 0, winnersNumber)

	iterations := 0
	for {
		if len(winners) > winnersNumber || len(participants) == 0 {
			break
		}

		if iterations > 100 {
			return []entity.Winner{}, fmt.Errorf("can`t generate random value while find winners in giveaway %d", giveawayID)
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
			winners = append(winners, newWinner(participants[index]))
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

func newWinner(participant entity.Participant) entity.Winner {
	return entity.Winner{
		ParticipantID: participant.ID,
		GiveawayID: participant.GiveawayID,
	}
}