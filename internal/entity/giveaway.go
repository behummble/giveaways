package entity

import (
	"time"

)

type Client struct {
	ID int
	CreatedAt time.Time
	UpdatetAt time.Time
	Name string
	UserID int
	AccessToken string
	RefreshToken string
}

type Giveaway struct {
	ID int
	CreatedAt time.Time
	Description string
	EventDate time.Time
	PublishDate time.Time
	WinnersNumber int
	ParticipantNumber int
	Terms []string
	ClientID int
	Client Client
	MessageID int
	GiveawayChatID int64
	SchedulerPublishID int
	SchedulerGiveawayID int
}

type Participant struct {
	ID int
	CreatedAt time.Time
	Ticket int64
	Username string
	UserID int64
	GiveawayID int
	Giveaway Giveaway
}

type Winner struct {
	ID int
	CreatedAt time.Time
	ParticipantID int
	Participant Participant
	GiveawayID int
	Giveaway Giveaway
}
