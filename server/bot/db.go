package bot

import (
	"app/models"
	"log"
)

func (s *Service) changeCountry(chatId int64, countryId int) {
	_, err := s.db.Exec("INSERT INTO user_countries(chatId, countryId) VALUES ($1, $2)", chatId, countryId)
	if err != nil {
		_, _ = s.db.Exec("UPDATE user_countries SET countryId = $1 WHERE chatId = $2", countryId, chatId)
	}
}

func (s *Service) getCountryByCode(code string) (*models.Country, error) {
	country := models.Country{}
	err := s.db.Get(&country, "SELECT * FROM countries WHERE code=$1", code)
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func (s *Service) getCountryIdByChatId(chatId int64) (int, error) {
	var countryId int
	err := s.db.Get(&countryId, "SELECT countryId FROM user_countries WHERE chatId=$1", chatId)
	if err != nil {
		return 0, err
	}
	return countryId, nil
}

func (s *Service) getCountryById(countryId int) (*models.Country, error) {
	country := models.Country{}
	err := s.db.Get(&country, "SELECT * FROM countries WHERE id=$1", countryId)
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func (s *Service) getCountries(offset int) ([]models.Country, int) {
	var countries []models.Country
	var count int
	err := s.db.Select(&countries, "SELECT * FROM countries ORDER BY id LIMIT 10 OFFSET $1", offset)
	if err != nil {
		log.Println(err)
	}
	err = s.db.Get(&count, "SELECT COUNT(*) FROM countries")
	if err != nil {
		log.Println(err)
	}
	return countries, count
}

func (s *Service) getAnswer(question string, countryId int) (string, error) {
	var answer string
	err := s.db.Get(&answer, "SELECT answer FROM qa WHERE country_id=$1 AND question=$2", countryId, question)
	if err != nil {
		return "", err
	}
	return answer, nil
}

func (s *Service) askQuestion(chatId int64, countryId int, question string) {
	if question != "" {
		_, _ = s.db.Exec("INSERT INTO user_questions(chat_id, country_id, question, status) VALUES ($1, $2, $3, $4)", chatId, countryId, question, 0)
	}
}

func (s *Service) setUserStatus(chatId int64, status string) {
	_, err := s.db.Exec("INSERT INTO users_bot_status(chat_id, status) VALUES ($1, $2)", chatId, status)
	if err != nil {
		_, _ = s.db.Exec("UPDATE users_bot_status SET status = $1 WHERE chat_id = $2", status, chatId)
	}
}

func (s *Service) getUserStatus(chatId int64) string {
	var status string
	err := s.db.Get(&status, "SELECT status FROM users_bot_status WHERE chat_id = $1", chatId)
	if err != nil {
		return "UNKNOWN"
	}
	return status
}