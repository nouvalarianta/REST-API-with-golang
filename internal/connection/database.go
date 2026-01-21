package connection

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func GetDatabase(databaseURL string) *sql.DB {
	db, err := sql.Open("postgres", databaseURL)

	if err != nil {
		log.Fatal("failed to open connection: ", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping connection: ", err.Error())
	}

	return db
}
