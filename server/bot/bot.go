package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var bot *tgbotapi.BotAPI

func Start() error {
	var err error
	bot, err = tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		return err
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Location != nil {
			handleLocation(update.Message.Chat.ID, update.Message.Location)
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "ask":
				handleAsk(update.Message.Chat.ID)
			default:
				handleUndefinedCommand(update.Message.Chat.ID)
			}
		}


	}
	return nil
}

func handleAsk(chatId int64) {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Оберіть країну вручну", "Оберіть країну вручну"),

		),
	)
	msg := tgbotapi.NewMessage(chatId, "Відправте свою геолокацію або")
	msg.ReplyMarkup = numericKeyboard
	_, _ = bot.Send(msg)
}

func handleUndefinedCommand(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "На жаль, не можу впізнати команду")
	_, _ = bot.Send(msg)
}

func handleLocation(chatId int64, location *tgbotapi.Location) {

}