package botservice

import (
	"context"
	"log/slog"
	"time"
	"fmt"
	"strings"

	"github.com/behummble/giveaway/internal/entity"
	"github.com/robfig/cron/v3"
	"gopkg.in/telebot.v3"
)

const (
	congratulationsText = "Поздравляем! Вы выиграли в розыгрыше!"
	giveawayResultText = "Розыгрыш проведен. Победителями стали:"
	cancelGiveaway = "К сожалению, розыгрыш был отменен."
)

var (
	inlineParticipate = telebot.InlineButton {
		Text: "Участвовать",
		InlineQuery: "/add_participate",
	}
)

type DB interface {
	Participant(ctx context.Context, id int) (entity.Participant, error)
	Giveaway(ctx context.Context, id int) (entity.Giveaway, error)
	UpdateGiveaway(ctx context.Context, giveaway entity.Giveaway) (entity.Giveaway, error)
	Winners(ctx context.Context, giveawayID int) ([]entity.Winner, error)
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

func (bot Bot) NotifyWinners(ctx context.Context, giveawayID int) {
	winners, err := bot.db.Winners(ctx, giveawayID)
	if err != nil {
		bot.log.Error(fmt.Sprintf("Can` get winners giveaway %d from db", giveawayID), "Error", err)
		return
	}

	giveaway, err := bot.db.Giveaway(ctx, giveawayID)
	if err != nil {
		bot.log.Error(fmt.Sprintf("Can`t get giveaway %d from db", giveawayID), "Error", err)
		return
	}

	for _, winner := range winners {
		go bot.notifyWinner(winner, giveaway.GiveawayChatID, giveaway.MessageID)
	}
}

func (bot Bot) PublishResults(ctx context.Context, giveawayID int) {
	giveaway, err := bot.db.Giveaway(ctx, giveawayID)
	if err != nil {
		bot.log.Error(fmt.Sprintf("Can`t get giveaway %d from db", giveawayID), "Error", err)
		return
	}

	winners, err := bot.db.Winners(ctx, giveawayID)
	if err != nil {
		bot.log.Error(fmt.Sprintf("Can` get winners giveaway %d from db", giveawayID), "Error", err)
		return
	}
	
	newText := publishText(winners)

	err = bot.editPublicationMessage(newText, giveaway.MessageID, giveaway.GiveawayChatID)
	if err != nil {
		bot.log.Error(
			fmt.Sprintf("Can` edit giveaway %d message %d in chat %d", 
				giveawayID, giveaway.MessageID, giveaway.GiveawayChatID), 
			"Error", err)
	}
}

func (bot Bot) CancelGiveaway(ctx context.Context, giveawayID int) {
	giveaway, err := bot.db.Giveaway(context.Background(), giveawayID)
	if err != nil {
		bot.log.Error(fmt.Sprintf("Can`t get giveaway %d from db", giveawayID), "Error", err)
		return
	}

	err = bot.cancelSchedulePublication(
		giveaway.PublishDate, 
		giveaway.MessageID,
		giveaway.GiveawayChatID,
		giveaway.SchedulerPublishID,
		giveawayID)
	
	if err != nil {
		bot.log.Error(
			fmt.Sprintf("Can` cancel schedule publication for giveaway %d message %d in chat %d", 
				giveawayID, giveaway.MessageID, giveaway.GiveawayChatID), 
			"Error", err)
	}	
}

func (bot Bot) SchedulePublication(ctx context.Context, publishDate time.Time, chatID int64, giveawayID int) (int, error) {
	specPublication := fmt.Sprintf(
		"%d %d %d %d %d ? %d",
		publishDate.Second(),
		publishDate.Minute(),
		publishDate.Hour(),
		publishDate.Day(),
		publishDate.Month(),
		publishDate.Year())
	idFunc, err := bot.cron.AddFunc(specPublication, schedulePublicationFunc(bot, giveawayID, chatID))
	
	return int(idFunc), err
}

func (bot Bot) ParticipantMetTerms(ctx context.Context, terms []string, userID int64) (bool, error) {
	metTerms := true
	for _, chat := range terms {
		metTerms, err := userMetTerm(chat, userID, bot.tgbot)
		if err != nil {
			return false, err
		}
		if !metTerms {
			break
		}
	}

	return metTerms, nil
}

func schedulePublicationFunc(bot Bot, giveawayID int, chatID int64) func() {
	return func() {
		giveaway, err := bot.db.Giveaway(context.Background(), giveawayID)
		if err != nil {
			bot.log.Error(fmt.Sprintf("Can`t get giveaway %d from db", giveawayID), "Error", err)
			return
		}

		chat, err := bot.tgbot.ChatByID(chatID)
		if err != nil {
			bot.log.Error(fmt.Sprintf("Can`t initialize giveaway tg_chat %d", chatID), "Error", err)
			return
		}
		keyboard := make([][]telebot.InlineButton, 1)
		keyboard[0] = append(keyboard[0], inlineParticipate)
		markup := &telebot.ReplyMarkup{
			InlineKeyboard: keyboard,
		}

		msg := &telebot.Message{
			Text: giveaway.Description,
			ReplyMarkup: markup,
		}

		tgMessage, err := bot.tgbot.Send(chat, msg)
		if err != nil {
			bot.log.Error(fmt.Sprintf("Can`t publish giveaway message in giveaway tg_chat %d", chatID), "Error", err)
			return
		}

		giveaway.MessageID = tgMessage.ID

		_, err = bot.db.UpdateGiveaway(context.Background(), giveaway)
		if err != nil {
			bot.log.Error(fmt.Sprintf("Can`t update giveaway %d in db", giveawayID), "Error", err)
			return
		}
	}
}

func (bot Bot) notifyWinner(winner entity.Winner, chatID int64, messageID int) {
	recipient, err := bot.tgbot.ChatByUsername(winner.Participant.Username)
	if err != nil {
		bot.log.Error(fmt.Sprintf("Can`t initialize user tg_chat %s", winner.Participant.Username), "Error", err)
		return
	}

	_, err = bot.tgbot.Send(recipient, congratulationsText)
	if err != nil {
		bot.log.Error(fmt.Sprintf("Can`t send message to user %s", winner.Participant.Username), "Error", err)
		return
	}

	giveawayChat, err := bot.tgbot.ChatByID(chatID)
	if err != nil {
		bot.log.Error(fmt.Sprintf("Can`t initialize giveaway tg_chat %d", chatID), "Error", err)
		return
	}

	giveawayMessage := &telebot.Message{
		ID: messageID,
		Chat: giveawayChat,
	}

	_, err = bot.tgbot.Forward(recipient, giveawayMessage)
	if err != nil {
		bot.log.Error(fmt.Sprintf("Can`t forward giveaway message to %s", winner.Participant.Username), "Error", err)
		return
	}

}

func publishText(winners []entity.Winner) string {
	var builder strings.Builder
	fmt.Fprint(&builder, giveawayResultText)
	for i, _ := range winners {
		fmt.Fprintln(&builder, winners[i].Participant.Username)
	}

	return builder.String()
}

func userMetTerm(chatName string, userID int64, tgBot *telebot.Bot) (bool, error) {
	chat, err := tgBot.ChatByUsername(chatName)
	if err != nil {
		return false, err
	}

	user := &telebot.User{
		ID: userID,
	}
	
	_, err = tgBot.ChatMemberOf(chat, user)
	return err == nil, nil
}

func (bot Bot) editPublicationMessage(newText string, messageID int, chatID int64) error {
	giveawayChat, err := bot.tgbot.ChatByID(chatID)
	if err != nil {
		return err
	}

	giveawayMessage := &telebot.Message{
		ID: messageID,
		Chat: giveawayChat,
	}

	_, err = bot.tgbot.Edit(giveawayMessage, fmt.Sprintf("%s \n %s", giveawayMessage.Text, newText))
	
	return err
}

func (bot Bot) cancelSchedulePublication(publishDate time.Time, messageID int, chatID int64, schedulePublishID, giveawayID int) error {
	if time.Now().After(publishDate) {
		err := bot.editPublicationMessage(cancelGiveaway, messageID, chatID)
		if err != nil {
			return err
		}
	} else {
		bot.cron.Remove(cron.EntryID(schedulePublishID))
		bot.log.Info(fmt.Sprintf("Schedule publication func for giveaway %d was removed", giveawayID))
	}
	
	return nil
}