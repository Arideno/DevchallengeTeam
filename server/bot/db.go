package bot

import (
	"log"
)

func (s *Service) changeCountry(chatId int64, countryId int) {
	_, err := s.db.Exec("INSERT INTO user_countries(chatId, countryId) VALUES ($1, $2)", chatId, countryId)
	if err != nil {
		_, _ = s.db.Exec("UPDATE user_countries SET countryId = $1 WHERE chatId = $2", countryId, chatId)
	}
}

func (s *Service) getCountryByCode(code string) (*country, error) {
	country := country{}
	err := s.db.Get(&country, "SELECT * FROM countries WHERE code=$1", code)
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func (s *Service) checkIfUserHasCountry(chatId int64) int {
	var countryId int
	err := s.db.Get(&countryId, "SELECT countryId FROM user_countries WHERE chatId=$1", chatId)
	if err != nil {
		return 0
	}
	return countryId
}

func (s *Service) getCountryById(countryId int) (*country, error){
	country := country{}
	err := s.db.Get(&country, "SELECT * FROM countries WHERE id=$1", countryId)
	if err != nil {
		return nil, err
	}
	return &country, nil
}

func (s *Service) getCountries(offset int) ([]country, int) {
	var countries []country
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
