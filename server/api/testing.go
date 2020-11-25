package api

import (
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestStore(t *testing.T) (*APIServer, sqlmock.Sqlmock) {
	t.Helper()

	s := &APIServer{}
	db, mock, _ := sqlmock.New()
	s.db = sqlx.NewDb(db, "pgx")

	return s, mock
}