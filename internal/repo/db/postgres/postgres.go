package postgres

import (
	"fmt"
	"log/slog"
	"context"

	"github.com/behummble/giveaway/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	log *slog.Logger
	conn *gorm.DB
}

func New(log *slog.Logger, username, password, dbName, host string, port int) (Postgres, error) {
	config := dbConfig(dbName, host, username, password, port)

	db, err := gorm.Open(postgres.New(config))
	if err != nil {
		return Postgres{}, err
	}

	return Postgres{
		log: log,
		conn: db,
	}, nil
}

func( client Postgres) AddGiveaway(ctx context.Context, data entity.Giveaway) (entity.Giveaway, error) {
	res := client.conn.WithContext(ctx).Create(&data)
	if res.Error != nil {
		return entity.Giveaway{}, res.Error
	}

	return data, res.Error
}

func (client Postgres) Giveaway(ctx context.Context, id int) (entity.Giveaway, error) {
	var giveaway entity.Giveaway
	res := client.conn.WithContext(ctx).First(&giveaway, id)

	return giveaway, res.Error
}

func (client Postgres) UpdateGiveaway(ctx context.Context, updGiveaway entity.Giveaway) (entity.Giveaway, error) {
	giveaway, err := client.Giveaway(ctx, updGiveaway.ID)
	if err != nil {
		return giveaway, err
	}

	res := client.conn.WithContext(ctx).Save(&updGiveaway)
	return updGiveaway, res.Error
}

func (client Postgres) DeleteGiveaway(ctx context.Context, id int) error {
	return client.conn.WithContext(ctx).Delete(&entity.Giveaway{}, id).Error
}

func (client Postgres) AddParticipant(ctx context.Context, data entity.Participant) (entity.Participant, error) {
	res := client.conn.WithContext(ctx).Create(&data)
	if res.Error != nil {
		return entity.Participant{}, res.Error
	}

	return data, res.Error
}

func (client Postgres) Participant(ctx context.Context, id int) (entity.Participant, error) {
	var participant entity.Participant
	res := client.conn.WithContext(ctx).First(&participant, id)

	return participant, res.Error
}

func (client Postgres) AllParticipants(ctx context.Context, giveawayID int) ([]entity.Participant, error) {
	var participants []entity.Participant
	res := client.conn.WithContext(ctx).Where("giveaway_id = ?", giveawayID).Find(&participants)

	return participants, res.Error
}

func (client Postgres) DeleteAllGiveawayParticipants(ctx context.Context, giveawayID int) error {
	return client.conn.WithContext(ctx).Where("giveaway_id = ?", giveawayID).Delete(&entity.Giveaway{}, giveawayID).Error
}

func (client Postgres) AddWinners(ctx context.Context, winners []entity.Winner) error {
	return client.conn.WithContext(ctx).Create(&winners).Error
}

func (client Postgres) Winners(ctx context.Context, giveawayID int) ([]entity.Winner, error) {
	var winners []entity.Winner
	res := client.conn.WithContext(ctx).Where("giveaway_id = ?", giveawayID).Find(&winners)
	
	return winners, res.Error
}

func (client Postgres) Client(ctx context.Context, key string) (entity.Client, error) {
	var customer entity.Client
	res := client.conn.WithContext(ctx).Where("api_key = ?", key).First(&customer)
	
	return customer, res.Error
}

func (client Postgres) AddClient(ctx context.Context, data entity.Client) (int, error) {
	res := client.conn.WithContext(ctx).Create(&data)
	if res.Error != nil {
		return -1, res.Error
	}

	return data.ID, res.Error
}

func (client Postgres) DeleteClient(ctx context.Context, id int) error {
	return client.conn.WithContext(ctx).Delete(&entity.Client{}, id).Error
}

func dbConfig(host, username, password, dbName string, port int) postgres.Config {
	return postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=enable TimeZone=Europe/Moscow",
						host, username, password, dbName, port),
	}
}	