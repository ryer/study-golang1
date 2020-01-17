package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

func Main(dbFile string) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = InsertData(db, []string{"flower", "book", "candy"})
	if err != nil {
		panic(err)
	}
}

// productsテーブルに指定したデータ（data）を挿入します。
func InsertData(db *sql.DB, data []string) (err error) {
	tx, e := db.Begin()
	if e != nil {
		err = errors.Wrap(e, "Failed to Begin")
		return
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			//noinspection GoUnhandledErrorResult
			tx.Rollback()
		}
	}()

	stmt, e := tx.Prepare(`INSERT INTO "products" ("name") VALUES (?)`)
	if e != nil {
		err = errors.Wrap(e, "Failed to Prepare")
		return
	}
	defer func() {
		e := stmt.Close()
		if e != nil {
			err = errors.Wrap(e, "Failed to Close")
		}
	}()

	for _, v := range data {
		_, e = stmt.Exec(v)
		if e != nil {
			err = errors.Wrap(e, "Failed to Exec")
			break
		}
	}

	return
}
