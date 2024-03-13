package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5433
	username = "postgres"
	password = "pgadmin"
	dbName   = "go_tiny_url_db"
)

func ConnectDB() *sql.DB {
	connectionUrl := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbName)
	db, err := sql.Open("postgres", connectionUrl)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
