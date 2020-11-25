package api

import (
	"app/models"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetQuestionsByCountryId(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test no questions", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM user_questions").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "chat_id", "country_id", "question", "status"})).RowsWillBeClosed()
		questions, err := s.getQuestionsByCountryId(1)
		assert.NoError(t, err)
		assert.Len(t, questions, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test one question", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM user_questions").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "chat_id", "country_id", "question", "status"}).FromCSVString("1,1,1,test,0"))
		questions, err := s.getQuestionsByCountryId(1)
		assert.NoError(t, err)
		assert.Len(t, questions, 1)
		assert.Equal(t, questions[0].Question, "test")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test sort ids", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM user_questions").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "chat_id", "country_id", "question", "status"}).FromCSVString("2,1,1,test2,1").FromCSVString("1,1,1,test,0"))
		questions, err := s.getQuestionsByCountryId(1)
		assert.NoError(t, err)
		assert.Len(t, questions, 2)
		assert.Equal(t, questions[0].Question, "test2")
		assert.Equal(t, questions[1].Question, "test")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetQuestionById(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test no questions", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM user_questions").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "chat_id", "country_id", "question", "status"})).RowsWillBeClosed()
		question, err := s.getQuestionById(1, 1)
		assert.Error(t, err)
		assert.Nil(t, question)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test one question", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM user_questions").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "chat_id", "country_id", "question", "status"}).FromCSVString("1,1,1,test,0"))
		question, err := s.getQuestionById(1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, question)
		assert.Equal(t, question.Question, "test")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSetStatus(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test error", func(t *testing.T) {
		mock.ExpectExec("UPDATE user_questions").WithArgs(0, 1, 1).WillReturnError(sql.ErrNoRows)
		err := s.setStatus(0, 1, 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test no error", func(t *testing.T) {
		mock.ExpectExec("UPDATE user_questions").WithArgs(0, 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
		err := s.setStatus(0, 1, 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAddQA(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	mock.ExpectExec("INSERT INTO qa").WithArgs(1, "test", "test", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	err := s.addQA("test", "test", 1, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetQAByCountry(t *testing.T)  {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test no qas", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "country_id", "topic_id", "question", "answer", "topic_name"})).RowsWillBeClosed()
		qas, err := s.getQAByCountry(1)
		assert.NoError(t, err)
		assert.Len(t, qas, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test one qas", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "country_id", "topic_id", "question", "answer", "topic_name"}).FromCSVString("1,1,1,test,test,test"))
		qas, err := s.getQAByCountry(1)
		assert.NoError(t, err)
		assert.Len(t, qas, 1)
		assert.Equal(t, qas[0].Question, "test")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test two qas", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "country_id", "topic_id", "question", "answer", "topic_name"}).FromCSVString("1,1,1,test,test,test").FromCSVString("2,1,1,test2,test,test"))
		qas, err := s.getQAByCountry(1)
		assert.NoError(t, err)
		assert.Len(t, qas, 2)
		assert.Equal(t, qas[0].Question, "test")
		assert.Equal(t, qas[1].Question, "test2")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetQAById(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test no qa", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "country_id", "topic_id", "question", "answer", "topic_name"})).RowsWillBeClosed()
		qa, err := s.getQAById(1, 1)
		assert.Error(t, err)
		assert.Nil(t, qa)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test one qa", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "country_id", "topic_id", "question", "answer", "topic_name"}).FromCSVString("1,1,1,test,test,test"))
		qa, err := s.getQAById(1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, qa)
		assert.Equal(t, qa.Question, "test")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateQA(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test error", func(t *testing.T) {
		mock.ExpectExec("UPDATE").WithArgs("test", "test", 1, 1, 1).WillReturnError(sql.ErrNoRows)
		err := s.updateQA(models.QA{
			Id:        1,
			CountryId: 1,
			TopicId:   1,
			Question:  "test",
			Answer:    "test",
		})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test no error", func(t *testing.T) {
		mock.ExpectExec("UPDATE").WithArgs("test", "test", 1, 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
		err := s.updateQA(models.QA{
			Id:        1,
			CountryId: 1,
			TopicId:   1,
			Question:  "test",
			Answer:    "test",
		})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteQA(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test error", func(t *testing.T) {
		mock.ExpectExec("DELETE").WithArgs(1, 1).WillReturnError(sql.ErrNoRows)
		err := s.deleteQA(models.QA{Id: 1, CountryId: 1})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test no error", func(t *testing.T) {
		mock.ExpectExec("DELETE").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
		err := s.deleteQA(models.QA{Id: 1, CountryId: 1})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUser(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs("test").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "country_id"})).RowsWillBeClosed()
		user, err := s.getUser("test")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test no error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs("test").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "country_id"}).FromCSVString("1,test,test,1"))
		user, err := s.getUser("test")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, user.Username, "test")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSendMessage(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	mock.ExpectQuery("INSERT").WithArgs(1, "test", 1, true).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
	message := s.sendMessage(1, 1, "test")
	assert.Equal(t, message.Message, "test")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMessages(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test not question", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "chat_id", "message", "from_operator", "question_id"})).RowsWillBeClosed()
		questions := s.getMessages(1)
		assert.Len(t, questions, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test one question", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "chat_id", "message", "from_operator", "question_id"}).FromCSVString("1,1,test,true,1"))
		questions := s.getMessages(1)
		assert.Len(t, questions, 1)
		assert.Equal(t, questions[0].Message, "test")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test two questions", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "chat_id", "message", "from_operator", "question_id"}).FromCSVString("1,1,test,true,1").FromCSVString("1,1,test2,false,1"))
		questions := s.getMessages(1)
		assert.Len(t, questions, 2)
		assert.Equal(t, questions[0].Message, "test")
		assert.Equal(t, questions[1].Message, "test2")
		assert.Equal(t, questions[0].FromOperator, true)
		assert.Equal(t, questions[1].FromOperator, false)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetChatIdByQuestionId(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"chat_id"}).FromCSVString("1"))
	chatId := s.getChatIdByQuestionId(1)
	assert.Equal(t, int64(1), chatId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTopics(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test no topics", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"})).RowsWillBeClosed()
		topics := s.getTopics()
		assert.Len(t, topics, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test one topic", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).FromCSVString("1,test"))
		topics := s.getTopics()
		assert.Len(t, topics, 1)
		assert.Equal(t, "test", topics[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test two topics", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).FromCSVString("1,test").FromCSVString("2,test2"))
		topics := s.getTopics()
		assert.Len(t, topics, 2)
		assert.Equal(t, "test", topics[0].Name)
		assert.Equal(t, "test2", topics[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestChangePassword(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("No user id", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"password"})).RowsWillBeClosed()
		err := s.changePassword(1, "test", "test2")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test passwords don't match", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"password"}).FromCSVString("test3"))
		err := s.changePassword(1, "test", "test2")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test ok", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"password"}).FromCSVString("$2a$04$YS1k0R.QaKJbF7U/UJdG/eY0tEm193vneUtVj1oOsg6ljUK5hiNS6"))
		err := s.changePassword(1, "test", "test2")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}