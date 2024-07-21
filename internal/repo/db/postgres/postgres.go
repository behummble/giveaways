package postgres

import (
	"fmt"
	"log/slog"

	"github.com/behummble/giveaway/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	log *slog.Logger
	conn *gorm.DB
}

func New(log *slog.Logger, dbName, host, username, password string, port int) Postgres {
	config := dbConfig(dbName, host, username, password, port)

	db, err := gorm.Open(postgres.New(config))
	if err != nil {
		panic(err)
	}

	return Postgres{
		log: log,
		conn: db,
	}
}

func(client Postgres) AddGiveaway(data entity.Giveaway) (int, error) {
	res := client.conn.Create(&data)
	if res.Error != nil {
		return -1, res.Error
	}

	return data.ID, res.Error
}

func(client Postgres) Giveaway(id int) (entity.Giveaway, error) {
	var giveaway entity.Giveaway
	res := client.conn.First(&giveaway, id)

	return giveaway, res.Error
}

func(client Postgres) UpdateGiveaway(updGiveaway entity.Giveaway) (entity.Giveaway, error) {
	giveaway, err := client.Giveaway(updGiveaway.ID)
	if err != nil {
		return giveaway, err
	}
	res := client.conn.Save(&updGiveaway)
	return updGiveaway, res.Error
}

func(client Postgres) DeleteGiveaway(id int) error {
	return client.conn.Delete(&entity.Giveaway{}, id).Error
}

func(client Postgres) Winners(id int) ([]entity.Winner, error) {
	var winners []entity.Winner
	res := client.conn.Where("giveaway_id = ?", id).Find(&winners)
	
	return winners, res.Error
}

func(client Postgres) Client(key string) (entity.Client, error) {
	var customer entity.Client
	res := client.conn.Where("api_key = ?", key).First(&customer)
	
	return customer, res.Error
}

func dbConfig(dbName, host, username, password string, port int) postgres.Config {
	return postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=enable TimeZone=Europe/Moscow",
						host, username, password, dbName, port),
	}
}	