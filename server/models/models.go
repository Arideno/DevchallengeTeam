package models

type Country struct {
	Id int
	Name string
	Code string
}

type Question struct {
	Id int `json:"id"`
	ChatId int64 `json:"chat_id" db:"chat_id"`
	CountryId int `json:"country_id" db:"country_id"`
	Question string `json:"question"`
	Status int `json:"status"`
}

type Topic struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type QA struct {
	Id int `json:"id"`
	CountryId int `json:"country_id" db:"country_id"`
	TopicId int `json:"topic_id" db:"topic_id"`
	Question string `json:"question"`
	Answer string `json:"answer"`
}

type User struct {
	Id int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	CountryId int `json:"country_id" db:"country_id"`
}

type Message struct {
	Id int `json:"id"`
	ChatId int64 `json:"chat_id" db:"chat_id"`
	QuestionId int `json:"question_id" db:"question_id"`
	Message string `json:"message"`
	FromOperator bool `json:"from_operator" db:"from_operator"`
}