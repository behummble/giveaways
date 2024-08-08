package bot

import (
	"log/slog"
	"gopkg.in/telebot.v3"
	"github.com/behummble/giveaway/internal/service/giveaway"
	//"github.com/behummble/giveaway/internal/middleware"
)

const (
	
)

type Bot struct {
	log *slog.Logger
	tgbot *telebot.Bot
	giveaway *giveawayservice.Giveaway
}

func New(log *slog.Logger, giveaway *giveawayservice.Giveaway, tgbot *telebot.Bot) (*Bot, error) {
	bot := &Bot{
		log: log,
		tgbot: tgbot,
		giveaway: giveaway,
	}

	return bot, nil
}

func (bot *Bot) Register() {
	adminsOnly := bot.tgbot.Group()
	registerAdminHandlers(adminsOnly, bot)
	registerCommonHandlers(bot)
	//bot.tgbot.Use(middleware.ClientExist())
}

func registerAdminHandlers(group *telebot.Group, bot *Bot) {
	//group.Use(middleware.IsAdmin())
	group.Handle("/list", bot.giveaways)
	group.Handle("/info", bot.giveaways)
	group.Handle("/new_giveaway", bot.giveaways)
	group.Handle("/delete_giveaway", bot.giveaways)
	group.Handle("/update_giveaway", bot.giveaways)
	group.Handle("/winners", bot.giveaways)
}

func registerCommonHandlers(bot *Bot) {
	bot.tgbot.Handle("/refresh_token", bot.giveaways)
	bot.tgbot.Handle("/register_client", bot.giveaways)
}

func (bot *Bot) giveaways(upd telebot.Context) error {
	return nil
}

func (bot *Bot) giveawayInfo(upd telebot.Context) error {

}

func (bot *Bot) registerClient(upd telebot.Context) error {
	msg := &telebot,Message{
		
	}
	return upd.Reply(`
		Укажите, пожалуйста, каналы, разделенные через ;, в 
		которых планируется проводить розыгрыши. 
		Введенное сообщение должно быть в формате 
		"@пример1; @пример2..."`)

}

func (bot *Bot) refreshToken(upd telebot.Context) error {
	
}

func (bot *Bot) newGiveaway(upd telebot.Context) error {
	
}

func (bot *Bot) updateGiveaway(upd telebot.Context) error {
	
}

func (bot *Bot) deleteGiveaway(upd telebot.Context) error {
	
}

func (bot *Bot) winners(upd telebot.Context) error {
	
}