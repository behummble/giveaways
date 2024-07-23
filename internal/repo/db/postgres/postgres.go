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

func(client Postgres) AddGiveaway(ctx context.Context, data entity.Giveaway) (int, error) {
	res := client.conn.WithContext(ctx).Create(&data)
	if res.Error != nil {
		return -1, res.Error
	}

	return data.ID, res.Error
}

func(client Postgres) Giveaway(ctx context.Context, id int) (entity.Giveaway, error) {
	var giveaway entity.Giveaway
	res := client.conn.WithContext(ctx).First(&giveaway, id)

	return giveaway, res.Error
}

func(client Postgres) UpdateGiveaway(ctx context.Context, updGiveaway entity.Giveaway) (entity.Giveaway, error) {
	giveaway, err := client.Giveaway(ctx, updGiveaway.ID)
	if err != nil {
		return giveaway, err
	}

	res := client.conn.WithContext(ctx).Save(&updGiveaway)
	return updGiveaway, res.Error
}

func(client Postgres) DeleteGiveaway(ctx context.Context, id int) error {
	return client.conn.WithContext(ctx).Delete(&entity.Giveaway{}, id).Error
}

func(client Postgres) AddParticipant(ctx context.Context, data entity.Participant) (int, error) {
	res := client.conn.WithContext(ctx).Create(&data)
	if res.Error != nil {
		return -1, res.Error
	}

	return data.ID, res.Error
}

func(client Postgres) Participant(ctx context.Context, id int) (entity.Participant, error) {
	var giveaway entity.Participant
	res := client.conn.WithContext(ctx).First(&giveaway, id)

	return giveaway, res.Error
}

func(client Postgres) Winners(ctx context.Context, id int) ([]entity.Winner, error) {
	var winners []entity.Winner
	res := client.conn.WithContext(ctx).Where("giveaway_id = ?", id).Find(&winners)
	
	return winners, res.Error
}

func(client Postgres) Client(ctx context.Context, key string) (entity.Client, error) {
	var customer entity.Client
	res := client.conn.WithContext(ctx).Where("api_key = ?", key).First(&customer)
	
	return customer, res.Error
}

func(client Postgres) AddClient(ctx context.Context, data entity.Client) (int, error) {
	res := client.conn.WithContext(ctx).Create(&data)
	if res.Error != nil {
		return -1, res.Error
	}

	return data.ID, res.Error
}

func(client Postgres) DeleteClient(ctx context.Context, id int) error {
	return client.conn.WithContext(ctx).Delete(&entity.Participant{}, id).Error
}

func dbConfig(host, username, password, dbName string, port int) postgres.Config {
	return postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=enable TimeZone=Europe/Moscow",
						host, username, password, dbName, port),
	}
}	