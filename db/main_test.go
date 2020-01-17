package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertData(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	pMock := mock.ExpectPrepare(`INSERT INTO "products" ("name") VALUES (?)`)
	pMock.ExpectExec().WithArgs("test1").WillReturnResult(sqlmock.NewResult(1, 1))
	pMock.ExpectExec().WithArgs("test2").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = InsertData(db, []string{"test1", "test2"})
	assert.NoError(t, err, "InsertData")
	assert.NoError(t, mock.ExpectationsWereMet(), "Valid SQL")
}
