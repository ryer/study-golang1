package gormdb

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertData(t *testing.T) {
	mockDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//noinspection GoUnhandledErrorResult
	defer mockDb.Close()

	db, err := gorm.Open("sqlite3", mockDb)
	assert.NoError(t, err, "gorm.Open")

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "products" ("name") VALUES (?)`).
		WithArgs("test1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO "products" ("name") VALUES (?)`).
		WithArgs("test2").
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	err = InsertData(db, []string{"test1", "test2"})
	assert.NoError(t, err, "InsertData")
	assert.NoError(t, mock.ExpectationsWereMet(), "Valid SQL")
}
