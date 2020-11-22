package bot

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