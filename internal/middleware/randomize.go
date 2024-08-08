package middleware

import (
	"gopkg.in/telebot.v3"
	"github.com/behummble/giveaway/internal/service/giveaway"

)

func IsAdmin(giveaway *giveawayservice.Giveaway, next telebot.HandlerFunc) telebot.HandlerFunc {
	return func (c telebot.Context) error {
		user := c.Message().Sender.ID
	}
}

func ClientExist() telebot.HandlerFunc {
	
}