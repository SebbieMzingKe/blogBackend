package database

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitializeDB(connStr string)  {
	DB, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Error connecting to the database.", err)
	}
	err = DB.Ping()

	if err != nil {
		log.Fatal("Database Unreachable")
	}
	log.Println("Database connected successfully")
}

