package bot

import (
	"fmt"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"strings"
)

type Service struct {
	bot *tgbotapi.BotAPI
	db *sqlx.DB
}

func (s *Service) Start() error {
	var err error
	s.bot, err = tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		return err
	}

	s.db, err = sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	s.bot.Debug = false

	log.Printf("Authorized on account %s", s.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Location != nil {
			s.handleLocation(update.Message.Chat.ID, update.Message.Location)
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "ask":
				s.handleAsk(update.Message.Chat.ID)
			default:
				s.handleUndefinedCommand(update.Message.Chat.ID)
			}
		}


	}
	return nil
}

func (s *Service) handleAsk(chatId int64) {
	id := s.checkIfUserHasCountry(chatId)
	if  id == 0 {
		msg := tgbotapi.NewMessage(chatId, "Відправте свою геолокацію або")
		msg.ReplyMarkup = createLocationInlineKeyboard()
		_, _ = s.bot.Send(msg)
	} else {
		country, _ := s.getCountryById(id)
		msgText := fmt.Sprintf("Вибрана країна - %s.", country.Name)
		msg := tgbotapi.NewMessage(chatId, msgText)
		msg.ReplyMarkup = createAskInlineKeyboard()
		_, _ = s.bot.Send(msg)
	}

}

func (s *Service) handleUndefinedCommand(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "На жаль, не можу впізнати команду")
	_, _ = s.bot.Send(msg)
}

func (s *Service) handleLocation(chatId int64, location *tgbotapi.Location) {
	geocoder := openstreetmap.Geocoder()
	address, err := geocoder.ReverseGeocode(location.Latitude, location.Longitude)
	if err != nil {
		msg := tgbotapi.NewMessage(chatId, "На жаль, ми не змогли розпізнати геопозицію, спробуйте відправити ще раз або")
		msg.ReplyMarkup = createLocationInlineKeyboard()
		_, _ = s.bot.Send(msg)
		return
	}
	countryCode := strings.ToLower(address.CountryCode)
	country, err := s.getCountryByCode(countryCode)
	if err != nil {
		msg := tgbotapi.NewMessage(chatId, "На жаль, посольства України немає у цій країні. Спробуйте ще раз або")
		msg.ReplyMarkup = createLocationInlineKeyboard()
		_, _ = s.bot.Send(msg)
		return
	}
	s.changeCountry(chatId, country.Id)
	msgText := fmt.Sprintf("Вибрана країна - %s.", country.Name)
	msg := tgbotapi.NewMessage(chatId, msgText)
	msg.ReplyMarkup = createAskInlineKeyboard()
	_, _ = s.bot.Send(msg)
}

func createLocationInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Оберіть країну вручну", "Оберіть країну вручну"),
		),
	)
}

func createAskInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Задати питання", "Задати питання"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Обрати іншу країну", "Обрати іншу країну"),
		),
	)
}