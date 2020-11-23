package api

import "app/models"

func (a *APIServer) getQuestionsByCountryId(countryId int) ([]models.Question, error) {
	questions := make([]models.Question, 0)
	err := a.db.Select(&questions, "SELECT * FROM user_questions WHERE country_id = $1", countryId)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (a *APIServer) getQuestionById(id int) (*models.Question, error) {
	question := &models.Question{}
	err := a.db.Get(question, "SELECT * FROM user_questions WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (a *APIServer) setStatus(status int, questionId int) error {
	_, err := a.db.Exec("UPDATE user_questions SET status = $1 WHERE id = $2", status, questionId)
	return err
}

func (a *APIServer) addQA(question string, answer string, countryId int) error {
	_, err := a.db.Exec("INSERT INTO qa(country_id, question, answer) VALUES ($1, $2, $3)", countryId, question, answer)
	return err
}

func (a *APIServer) getQAByCountry(countryId int) ([]models.QA, error) {
	qa := make([]models.QA, 0)
	err := a.db.Select(&qa, "SELECT * FROM qa WHERE country_id = $1", countryId)
	if err != nil {
		return nil, err
	}
	return qa, nil
}

func (a *APIServer) getQAById(id int) (*models.QA, error) {
	qa := &models.QA{}
	err := a.db.Get(qa, "SELECT * FROM qa WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return qa, nil
}

func (a *APIServer) updateQA(qa models.QA) error {
	_, err := a.db.Exec("UPDATE qa SET question=$1, answer=$2 WHERE id = $3", qa.Question, qa.Answer, qa.Id)
	return err
}