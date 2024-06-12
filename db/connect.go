package db

import (
	"database/sql"
	"log"
)

var db *sql.DB

func InitDB(dataSourceName string) *sql.DB {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	return db
}
