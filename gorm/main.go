package gormdb

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Product struct {
	ID   uint `gorm:"primary_key"`
	Name string
}

func Main(dbFile string) {
	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	err = InsertData(db, []string{"flower", "book", "candy"})
	if err != nil {
		panic(err)
	}
}

// productsテーブルに指定したデータ（data）を挿入します。
func InsertData(db *gorm.DB, data []string) (err error) {
	tx := db.Begin()
	if tx.Error != nil {
		err = errors.Wrap(tx.Error, "Failed to Begin")
		return
	}

	defer func() {
		r := recover()
		if r != nil || err != nil {
			//noinspection GoUnhandledErrorResult
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	for _, v := range data {
		p := &Product{Name: v}
		e := tx.Create(p).Error
		if e != nil {
			err = errors.Wrap(e, "Failed to Create")
			break
		}
	}

	return
}
