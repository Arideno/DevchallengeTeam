package bot

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChangeCountry(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test chatId does not exist", func(t *testing.T) {
		mock.ExpectExec("INSERT").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
		s.changeCountry(1, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test chatId does exist", func(t *testing.T) {
		mock.ExpectExec("INSERT").WithArgs(1, 1).WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectExec("UPDATE").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
		s.changeCountry(1, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetCountryByCode(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs("de").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code"})).RowsWillBeClosed()
		country, err := s.getCountryByCode("de")
		assert.Error(t, err)
		assert.Nil(t, country)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test no error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs("de").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code"}).FromCSVString("1,Germany,de"))
		country, err := s.getCountryByCode("de")
		assert.NoError(t, err)
		assert.NotNil(t, country)
		assert.Equal(t, "Germany", country.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetCountryIdByChatId(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"countryId"})).RowsWillBeClosed()
		countryId, err := s.getCountryIdByChatId(1)
		assert.Error(t, err)
		assert.Equal(t, 0, countryId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test no error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"countryId"}).FromCSVString("1"))
		countryId, err := s.getCountryIdByChatId(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, countryId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetCountryById(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code"})).RowsWillBeClosed()
		country, err := s.getCountryById(1)
		assert.Error(t, err)
		assert.Nil(t, country)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test no error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code"}).FromCSVString("1,Germany,de"))
		country, err := s.getCountryById(1)
		assert.NoError(t, err)
		assert.NotNil(t, country)
		assert.Equal(t, "Germany", country.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetCountries(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test no countries", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code"})).RowsWillBeClosed()
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count"}).FromCSVString("0"))
		countries, count := s.getCountries(0)
		assert.Equal(t, 0, count)
		assert.Len(t, countries, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test one countries", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code"}).FromCSVString("1,Germany,de"))
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count"}).FromCSVString("1"))
		countries, count := s.getCountries(0)
		assert.Equal(t, 1, count)
		assert.Len(t, countries, 1)
		assert.Equal(t, "Germany", countries[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test two countries", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code"}).FromCSVString("1,Germany,de").FromCSVString("2,Ukraine,ua"))
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count"}).FromCSVString("2"))
		countries, count := s.getCountries(0)
		assert.Equal(t, 2, count)
		assert.Len(t, countries, 2)
		assert.Equal(t, "Germany", countries[0].Name)
		assert.Equal(t, "Ukraine", countries[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test offset", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code"}).FromCSVString("1,Germany,de"))
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count"}).FromCSVString("2"))
		countries, count := s.getCountries(1)
		assert.Equal(t, 2, count)
		assert.Len(t, countries, 1)
		assert.Equal(t, "Germany", countries[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetAnswer(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"topicId"}).FromCSVString("1"))
		mock.ExpectQuery("SELECT").WithArgs(1, "test", 1).WillReturnRows(sqlmock.NewRows([]string{"answer"})).RowsWillBeClosed()

		answer, err := s.getAnswer("test", 1, 1)
		assert.Error(t, err)
		assert.Equal(t, "", answer)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test no error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"topicId"}).FromCSVString("1"))
		mock.ExpectQuery("SELECT").WithArgs(1, "test", 1).WillReturnRows(sqlmock.NewRows([]string{"answer"}).FromCSVString("test"))

		answer, err := s.getAnswer("test", 1, 1)
		assert.NoError(t, err)
		assert.Equal(t, "test", answer)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAskQuestion(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	mock.ExpectExec("INSERT").WithArgs(1, 1, "test", 0).WillReturnResult(sqlmock.NewResult(1, 1))
	s.askQuestion(1, 1, "test")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetUserStatus(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test user does not exist", func(t *testing.T) {
		mock.ExpectExec("INSERT").WithArgs(1, "UNKNOWN").WillReturnResult(sqlmock.NewResult(1, 1))
		s.SetUserStatus(1, "UNKNOWN")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test user exists", func(t *testing.T) {
		mock.ExpectExec("INSERT").WithArgs(1, "UNKNOWN").WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectExec("UPDATE").WithArgs("UNKNOWN", 1).WillReturnResult(sqlmock.NewResult(1, 1))
		s.SetUserStatus(1, "UNKNOWN")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUserStatus(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test no status", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"status"})).RowsWillBeClosed()
		status := s.getUserStatus(1)
		assert.Equal(t, "UNKNOWN", status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test ok", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"status"}).FromCSVString("DISCUSS"))
		status := s.getUserStatus(1)
		assert.Equal(t, "DISCUSS", status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetTopicsList(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	t.Run("Test error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"})).WillReturnError(sqlmock.ErrCancelled)
		topics, err := s.getTopicsList()
		assert.Error(t, err)
		assert.Len(t, topics, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test no topics", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"})).RowsWillBeClosed()
		topics, err := s.getTopicsList()
		assert.NoError(t, err)
		assert.Len(t, topics, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test one topic", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).FromCSVString("1,test"))
		topics, err := s.getTopicsList()
		assert.NoError(t, err)
		assert.Len(t, topics, 1)
		assert.Equal(t, "test", topics[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test two topics", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).FromCSVString("1,test").FromCSVString("2,test2"))
		topics, err := s.getTopicsList()
		assert.NoError(t, err)
		assert.Len(t, topics, 2)
		assert.Equal(t, "test", topics[0].Name)
		assert.Equal(t, "test2", topics[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSetUserTopic(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	mock.ExpectExec("UPDATE").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	s.setUserTopic(1, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTopicNameByTopicId(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"name"}).FromCSVString("test"))
	name := s.getTopicNameByTopicId(1)
	assert.Equal(t, "test", name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLastQuestionId(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
	id := s.getLastQuestionId(1)
	assert.Equal(t, 1, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSendMessage(t *testing.T) {
	s, mock := TestStore(t)
	defer s.db.Close()

	mock.ExpectQuery("INSERT").WithArgs(1, "test", 1, false).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))
	message := s.sendMessage(1, 1, "test")
	assert.Equal(t, "test", message.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}