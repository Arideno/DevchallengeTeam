package api

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) (*APIServer, sqlmock.Sqlmock) {
	t.Helper()

	s := &APIServer{}
	db, mock, _ := sqlmock.New()
	s.db = sqlx.NewDb(db, "pgx")

	return s, mock
}

func TestToken(t *testing.T) string {
	t.Helper()

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
	_ = json.Unmarshal(responseBody, &res)
	return res.Token
}