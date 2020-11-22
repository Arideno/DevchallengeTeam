package bot

import (
	"encoding/json"
	"fmt"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"math"
	"os"
	"strings"
)

type Service struct {
	bot *tgbotapi.BotAPI
	db *sqlx.DB
}

type CallbackData struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
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
		if update.Message != nil {
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

		if update.CallbackQuery != nil {
			callbackData := CallbackData{}
			_ = json.Unmarshal([]byte(update.CallbackQuery.Data), &callbackData)
			if callbackData.Type == "country" {
				s.handleCountryCallback(update.CallbackQuery.Message.Chat.ID, 1)
			} else if callbackData.Type == "page" {
				page := int(callbackData.Data.(float64))
				_, _ = s.bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
				s.handleCountryCallback(update.CallbackQuery.Message.Chat.ID, page)
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

func (s *Service) handleCountryCallback(chatId int64, page int) {
	countries, count := s.getCountries((page - 1) * 10)
	var row []tgbotapi.InlineKeyboardButton
	if page > 1 {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("<", fmt.Sprintf(`{"type": "page", "data": %d}`, page-1)))
	}
	if page+1 <= int(math.Ceil(float64(count) / 10)) {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(">", fmt.Sprintf(`{"type": "page", "data": %d}`, page+1)))
	}
	markup := tgbotapi.NewInlineKeyboardMarkup(row)
	msgText := "Виберіть країну зі списку:\n"
	for i, country := range countries {
		msgText += fmt.Sprintf("%s - %s", country.Name, "/country_" + country.Code)
		if i != len(countries) - 1 {
			msgText += "\n"
		}
	}
	msg := tgbotapi.NewMessage(chatId, msgText)
	msg.ReplyMarkup = markup
	_, _ = s.bot.Send(msg)
}

func createLocationInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Оберіть країну вручну", `{"type": "country"}`),
		),
	)
}

func createAskInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Задати питання", "Задати питання"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Обрати іншу країну", `{"type": "country"}`),
		),
	)
}