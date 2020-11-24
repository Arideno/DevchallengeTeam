package api

import (
	"app/models"
	"log"
)

func (a *APIServer) getQuestionsByCountryId(countryId int) ([]models.Question, error) {
	questions := make([]models.Question, 0)
	err := a.db.Select(&questions, "SELECT * FROM user_questions WHERE country_id = $1 ORDER BY id DESC", countryId)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (a *APIServer) getQuestionById(id int, countryId int) (*models.Question, error) {
	question := &models.Question{}
	err := a.db.Get(question, "SELECT * FROM user_questions WHERE id = $1 AND country_id = $2", id, countryId)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (a *APIServer) setStatus(status int, questionId int, countryId int) error {
	_, err := a.db.Exec("UPDATE user_questions SET status = $1 WHERE id = $2 AND country_id = $3", status, questionId, countryId)
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

func (a *APIServer) getUser(username string) (*models.User, error) {
	user := &models.User{}
	err := a.db.Get(user, "SELECT * FROM users WHERE username=$1", username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *APIServer) createUser(user *models.User) (int, error) {
	var id int
	err := a.db.Get(&id, "INSERT INTO users(username, password, country_id) VALUES ($1, $2, $3) RETURNING id", user.Username, user.Password, user.CountryId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *APIServer) sendMessage(chatId int64, questionId int, message string) models.Message {
	var id int
	err := a.db.Get(&id, "INSERT INTO user_messages(chat_id, message, question_id, from_operator) VALUES ($1, $2, $3, $4) RETURNING id", chatId, message, questionId, true)
	if err != nil {
		log.Println(err)
	}
	return models.Message{
		Id:           id,
		ChatId:       chatId,
		QuestionId:   questionId,
		Message:      message,
		FromOperator: true,
	}
}

func (a *APIServer) getMessages(questionId int) []models.Message {
	messages := make([]models.Message, 0)
	_ = a.db.Select(&messages, "SELECT * FROM user_messages WHERE question_id = $1", questionId)
	return messages
}