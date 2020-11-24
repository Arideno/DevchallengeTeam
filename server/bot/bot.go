package bot

import (
	"app/models"
	"encoding/json"
	"fmt"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/websocket"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"math"
	"os"
	"strings"
)

type API interface {
	GetConnections() map[int]*websocket.Conn
}

type Service struct {
	bot *tgbotapi.BotAPI
	db  *sqlx.DB
	ApiServer API
}

type CallbackData struct {
	Type string      `json:"type"`
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
			if update.Message.IsCommand() {
				if s.getUserStatus(update.Message.Chat.ID) == "DISCUSS" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Будь ласка очікуйте.")
					_, _ = s.bot.Send(msg)
					continue
				}
				switch update.Message.Command() {
				case "ask":
					s.handleAsk(update.Message.Chat.ID)
				default:
					isCountryListCommand, code := checkIfCountryByListCommand(update.Message.Command())
					if isCountryListCommand {
						s.handleCountryFromList(update.Message.Chat.ID, code)
					} else {
						s.handleUndefinedCommand(update.Message.Chat.ID)
					}
				}
			} else if update.Message.Location != nil {
				if s.getUserStatus(update.Message.Chat.ID) == "DISCUSS" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Будь ласка очікуйте.")
					_, _ = s.bot.Send(msg)
					continue
				}
				s.handleLocation(update.Message.Chat.ID, update.Message.Location)
			} else if update.Message.Text != "" {
				switch s.getUserStatus(update.Message.Chat.ID) {
				case "QUESTION":
					s.getAnswerOnQuestion(update.Message.Chat.ID, update.Message.Text)
				case "DISCUSS":
					s.sendQuestionToOperator(update.Message.Chat.ID, update.Message.Text)
				default:
					log.Println("Unknown text")
				}
			}

		}

		if update.CallbackQuery != nil {
			if s.getUserStatus(update.CallbackQuery.Message.Chat.ID) == "DISCUSS" {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Будь ласка очікуйте.")
				_, _ = s.bot.Send(msg)
				continue
			}
			callbackData := CallbackData{}
			_ = json.Unmarshal([]byte(update.CallbackQuery.Data), &callbackData)
			if callbackData.Type == "country" {
				s.handleCountryCallback(update.CallbackQuery.Message.Chat.ID, 1)
			} else if callbackData.Type == "page" {
				page := int(callbackData.Data.(float64))
				_, _ = s.bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
				s.handleCountryCallback(update.CallbackQuery.Message.Chat.ID, page)
			} else if callbackData.Type == "question" {
				s.handleQuestion(update.CallbackQuery.Message.Chat.ID)
			} else if callbackData.Type == "topic" {
				_, _ = s.bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
				s.handleTopic(update.CallbackQuery.Message.Chat.ID)
			} else if callbackData.Type == "topicChoice" {
				_, _ = s.bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
				topicId := callbackData.Data.(float64)
				s.handleSetUserTopic(update.CallbackQuery.Message.Chat.ID, topicId)
			}
		}


	}
	return nil
}

func (s *Service) handleAsk(chatId int64) {
	id, err := s.getCountryIdByChatId(chatId)
	if err != nil {
		msg := tgbotapi.NewMessage(chatId, "Відправте свою геолокацію або")
		msg.ReplyMarkup = createLocationInlineKeyboard()
		_, _ = s.bot.Send(msg)
	} else {
		country, _ := s.getCountryById(id)
		msgText := fmt.Sprintf("Вибрана країна - %s.", country.Name)
		msg := tgbotapi.NewMessage(chatId, msgText)
		msg.ReplyMarkup = createTopicInlineKeyboard()
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
	msg.ReplyMarkup = createTopicInlineKeyboard()
	_, _ = s.bot.Send(msg)
}

func (s *Service) handleCountryCallback(chatId int64, page int) {
	countries, count := s.getCountries((page - 1) * 10)
	var row []tgbotapi.InlineKeyboardButton
	if page > 1 {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData("<", fmt.Sprintf(`{"type": "page", "data": %d}`, page-1)))
	}
	if page+1 <= int(math.Ceil(float64(count)/10)) {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(">", fmt.Sprintf(`{"type": "page", "data": %d}`, page+1)))
	}
	markup := tgbotapi.NewInlineKeyboardMarkup(row)
	msgText := "Виберіть країну зі списку:\n"
	for i, country := range countries {
		msgText += fmt.Sprintf("%s - %s", country.Name, "/country_"+country.Code)
		if i != len(countries)-1 {
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
			tgbotapi.NewInlineKeyboardButtonData("Задати питання", `{"type": "question"}`),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Обрати іншу тему", `{"type": "topic"}`),
		),
	)
}


func createTopicInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Обрати тему запитання", `{"type": "topic"}`),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Обрати іншу країну", `{"type": "country"}`),
		),
	)
}

func checkIfCountryByListCommand(command string) (bool, string) {
	splitCommand := strings.Split(command, "_")
	if len(splitCommand) == 2 && splitCommand[0] == "country" {
		return true, splitCommand[1]
	}
	return false, ""
}

func (s *Service) handleCountryFromList(chatId int64, code string) {
	country, err := s.getCountryByCode(code)
	if err != nil {
		s.handleUndefinedCommand(chatId)
		return
	}
	s.changeCountry(chatId, country.Id)
	msgText := fmt.Sprintf("Вибрана країна - %s.", country.Name)
	msg := tgbotapi.NewMessage(chatId, msgText)
	msg.ReplyMarkup = createTopicInlineKeyboard()
	_, _ = s.bot.Send(msg)

}

func (s *Service) handleQuestion(chatId int64) {
	s.SetUserStatus(chatId, "QUESTION")
	msg := tgbotapi.NewMessage(chatId, "Ми чекаємо на Ваше запитання.")
	_, _ = s.bot.Send(msg)
}


func (s *Service) handleTopic(chatId int64) {
	_, err := s.getCountryIdByChatId(chatId) // countryId instead _
	log.Println(err)
	if err != nil {
		msg := tgbotapi.NewMessage(chatId, "Відправте свою геолокацію або")
		msg.ReplyMarkup = createLocationInlineKeyboard()
		_, _ = s.bot.Send(msg)
		return
	}
	topics, _:= s.getTopicsList()
	msg := tgbotapi.NewMessage(chatId, "Оберіть одну із тем")
	msg.ReplyMarkup = s.getTopicsInlineKeyboard(topics)
	_, _ = s.bot.Send(msg)
	return


}

func (s *Service) getTopicsInlineKeyboard(topics []models.Topic) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for i := 0; i < len(topics); i++ {
		newRow := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(topics[i].Name, fmt.Sprintf(`{"type": "topicChoice", "data": %d}`, topics[i].Id)))
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, newRow)
	}
	return keyboard
}


func (s *Service) getAnswerOnQuestion(chatId int64, question string)  {
	countryId, err := s.getCountryIdByChatId(chatId)
	if err != nil {
		msg := tgbotapi.NewMessage(chatId, "Відправте свою геолокацію або")
		msg.ReplyMarkup = createLocationInlineKeyboard()
		_, _ = s.bot.Send(msg)
		return
	}
	answer, err := s.getAnswer(question, countryId, chatId)
	if err != nil {
		s.askQuestion(chatId, countryId, question)
		s.SetUserStatus(chatId, "DISCUSS")
		msg := tgbotapi.NewMessage(chatId, "Ваше питання було передано оператору. Будь ласка очікуйте на відповідь.")
		_, _ = s.bot.Send(msg)
		return
	}
	s.SetUserStatus(chatId, "UNKNOWN")
	msg := tgbotapi.NewMessage(chatId, answer)
	_, _ = s.bot.Send(msg)
}


func (s *Service) handleSetUserTopic(chatId int64, topicId float64) {
	s.setUserTopic(chatId, topicId)
	topicName := s.getTopicNameByTopicId(topicId)
	msgText := fmt.Sprintf("Вибрана тема - %s.", topicName)
	msg := tgbotapi.NewMessage(chatId, msgText)
	msg.ReplyMarkup = createAskInlineKeyboard()
	_, _ = s.bot.Send(msg)
}

func (s *Service) sendQuestionToOperator(chatId int64, message string) {
	questionId := s.getLastQuestionId(chatId)
	msg := s.sendMessage(chatId, questionId, message)

	if conn, ok := s.ApiServer.GetConnections()[questionId]; ok {
		_ = conn.WriteJSON(map[string]interface{}{
			"type": "message",
			"data": msg,
		})
	}
}

func (s *Service) SendMessage(chatId int64, message string) {
	msg := tgbotapi.NewMessage(chatId, message)
	_, _ = s.bot.Send(msg)
}