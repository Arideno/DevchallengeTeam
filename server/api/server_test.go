package api

import (
	"app/models"
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func TestAPIServerLoginHandler(t *testing.T) {
	t.Run("Test user not found", func(t *testing.T) {
		s, mock := TestServer(t)
		defer s.db.Close()
		s.r = gin.New()
		s.configureRouter()

		mock.ExpectQuery("SELECT").WithArgs("test").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "country_id"})).RowsWillBeClosed()
		request := gin.H{
			"username": "test",
			"password": "test",
		}
		body, _ := json.Marshal(request)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
		s.r.ServeHTTP(w, r)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test passwords don't match", func(t *testing.T) {
		s, mock := TestServer(t)
		defer s.db.Close()
		s.r = gin.New()
		s.configureRouter()

		mock.ExpectQuery("SELECT").WithArgs("test").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "country_id"}).FromCSVString("1,test,test,1"))
		request := gin.H{
			"username": "test",
			"password": "test",
		}
		body, _ := json.Marshal(request)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
		s.r.ServeHTTP(w, r)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Test ok", func(t *testing.T) {
		s, mock := TestServer(t)
		defer s.db.Close()
		s.r = gin.New()
		s.configureRouter()

		mock.ExpectQuery("SELECT").WithArgs("test").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "country_id"}).FromCSVString("1,test,$2a$04$YS1k0R.QaKJbF7U/UJdG/eY0tEm193vneUtVj1oOsg6ljUK5hiNS6,1"))
		request := gin.H{
			"username": "test",
			"password": "test",
		}
		body, _ := json.Marshal(request)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
		s.r.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		responseBody, err := ioutil.ReadAll(w.Result().Body)
		assert.NoError(t, err)
		defer w.Result().Body.Close()
		type response struct {
			Token string `json:"token" binding:"required"`
		}
		var res response
		err = json.Unmarshal(responseBody, &res)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAPIServer_HandleQuestions(t *testing.T) {
	token := TestToken(t)
	s, mock := TestServer(t)
	defer s.db.Close()
	s.r = gin.New()
	s.configureRouter()
	mock.ExpectQuery("SELECT \\* FROM user_questions").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "chat_id", "country_id", "question", "status"}).FromCSVString("1,1,1,test,0"))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/api/auth/questions", bytes.NewReader([]byte{}))
	r.Header.Set("Authorization", "Bearer " + token)
	s.r.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	defer w.Result().Body.Close()
	questions := make([]models.Question, 0)
	err = json.Unmarshal(responseBody, &questions)
	assert.NoError(t, err)
	assert.Len(t, questions, 1)
	assert.Equal(t, "test", questions[0].Question)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAPIServer_HandleGetQuestionById(t *testing.T) {
	token := TestToken(t)
	s, mock := TestServer(t)
	defer s.db.Close()
	s.r = gin.New()
	s.configureRouter()
	mock.ExpectQuery("SELECT \\* FROM user_questions").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "chat_id", "country_id", "question", "status"}).FromCSVString("1,1,1,test,0"))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/api/auth/questions/1", bytes.NewReader([]byte{}))
	r.Header.Set("Authorization", "Bearer " + token)
	s.r.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	defer w.Result().Body.Close()
	var question models.Question
	err = json.Unmarshal(responseBody, &question)
	assert.NoError(t, err)
	assert.Equal(t, "test", question.Question)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAPIServer_HandleChangeStatus(t *testing.T) {
	token := TestToken(t)
	s, mock := TestServer(t)
	defer s.db.Close()
	s.r = gin.New()
	s.configureRouter()
	mock.ExpectExec("UPDATE user_questions").WithArgs(0, 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	w := httptest.NewRecorder()
	request := gin.H{
		"status": 0,
		"question_id": 1,
	}
	body, _ := json.Marshal(request)
	r, _ := http.NewRequest(http.MethodPatch, "/api/auth/changeStatus", bytes.NewReader(body))
	r.Header.Set("Authorization", "Bearer " + token)
	s.r.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	defer w.Result().Body.Close()
	type response struct {
		Message string `json:"message"`
	}
	var res response
	err = json.Unmarshal(responseBody, &res)
	assert.NoError(t, err)
	assert.Equal(t, "ok", res.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAPIServer_HandleQAs(t *testing.T) {
	token := TestToken(t)
	s, mock := TestServer(t)
	defer s.db.Close()
	s.r = gin.New()
	s.configureRouter()
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "country_id", "topic_id", "question", "answer", "topic_name"}).FromCSVString("1,1,1,test,test,test"))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/api/auth/qa", bytes.NewReader([]byte{}))
	r.Header.Set("Authorization", "Bearer " + token)
	s.r.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	defer w.Result().Body.Close()
	qas := make([]models.QA, 0)
	err = json.Unmarshal(responseBody, &qas)
	assert.NoError(t, err)
	assert.Len(t, qas, 1)
	assert.Equal(t, "test", qas[0].Question)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAPIServer_HandleAddQA(t *testing.T) {
	token := TestToken(t)
	s, mock := TestServer(t)
	defer s.db.Close()
	s.r = gin.New()
	s.configureRouter()
	mock.ExpectExec("INSERT INTO qa").WithArgs(1, "test", "test", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	w := httptest.NewRecorder()
	request := gin.H{
		"topic_id": 1,
		"question": "test",
		"answer": "test",
	}
	body, _ := json.Marshal(request)
	r, _ := http.NewRequest(http.MethodPost, "/api/auth/qa", bytes.NewReader(body))
	r.Header.Set("Authorization", "Bearer " + token)
	s.r.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	defer w.Result().Body.Close()
	type response struct {
		Message string `json:"message"`
	}
	var res response
	err = json.Unmarshal(responseBody, &res)
	assert.NoError(t, err)
	assert.Equal(t, "ok", res.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAPIServer_HandleGetQAById(t *testing.T) {
	token := TestToken(t)
	s, mock := TestServer(t)
	defer s.db.Close()
	s.r = gin.New()
	s.configureRouter()
	mock.ExpectQuery("SELECT").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "country_id", "topic_id", "question", "answer", "topic_name"}).FromCSVString("1,1,1,test,test,test"))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/api/auth/qa/1", bytes.NewReader([]byte{}))
	r.Header.Set("Authorization", "Bearer " + token)
	s.r.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	defer w.Result().Body.Close()
	var qa models.QA
	err = json.Unmarshal(responseBody, &qa)
	assert.NoError(t, err)
	assert.Equal(t, "test", qa.Question)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAPIServer_HandleUpdateQAById(t *testing.T) {
	token := TestToken(t)
	s, mock := TestServer(t)
	defer s.db.Close()
	s.r = gin.New()
	s.configureRouter()
	mock.ExpectExec("UPDATE").WithArgs("test", "test", 1, 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	w := httptest.NewRecorder()
	request := gin.H{
		"topic_id": 1,
		"question": "test",
		"answer": "test",
	}
	body, _ := json.Marshal(request)
	r, _ := http.NewRequest(http.MethodPut, "/api/auth/qa/1", bytes.NewReader(body))
	r.Header.Set("Authorization", "Bearer " + token)
	s.r.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	defer w.Result().Body.Close()
	type response struct {
		Message string `json:"message"`
	}
	var res response
	err = json.Unmarshal(responseBody, &res)
	assert.NoError(t, err)
	assert.Equal(t, "ok", res.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAPIServer_HandleDeleteQAById(t *testing.T) {
	token := TestToken(t)
	s, mock := TestServer(t)
	defer s.db.Close()
	s.r = gin.New()
	s.configureRouter()
	mock.ExpectExec("DELETE").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodDelete, "/api/auth/qa/1", bytes.NewReader([]byte{}))
	r.Header.Set("Authorization", "Bearer " + token)
	s.r.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	defer w.Result().Body.Close()
	type response struct {
		Message string `json:"message"`
	}
	var res response
	err = json.Unmarshal(responseBody, &res)
	assert.NoError(t, err)
	assert.Equal(t, "ok", res.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAPIServer_HandleGetTopics(t *testing.T) {
	token := TestToken(t)
	s, mock := TestServer(t)
	defer s.db.Close()
	s.r = gin.New()
	s.configureRouter()
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).FromCSVString("1,test"))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/api/auth/topics", bytes.NewReader([]byte{}))
	r.Header.Set("Authorization", "Bearer " + token)
	s.r.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	defer w.Result().Body.Close()
	topics := make([]models.Topic, 0)
	err = json.Unmarshal(responseBody, &topics)
	assert.NoError(t, err)
	assert.Len(t, topics, 1)
	assert.Equal(t, "test", topics[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAPIServer_HandleChangePassword(t *testing.T) {
	token := TestToken(t)
	s, mock := TestServer(t)
	defer s.db.Close()
	s.r = gin.New()
	s.configureRouter()
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"password"}).FromCSVString("$2a$04$YS1k0R.QaKJbF7U/UJdG/eY0tEm193vneUtVj1oOsg6ljUK5hiNS6"))
	w := httptest.NewRecorder()
	request := gin.H{
		"current_password": "test",
		"new_password": "test2",
	}
	body, _ := json.Marshal(request)
	r, _ := http.NewRequest(http.MethodPatch, "/api/auth/change/password", bytes.NewReader(body))
	r.Header.Set("Authorization", "Bearer " + token)
	s.r.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	defer w.Result().Body.Close()
	type response struct {
		Message string `json:"message"`
	}
	var res response
	err = json.Unmarshal(responseBody, &res)
	assert.NoError(t, err)
	assert.Equal(t, "ok", res.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}