package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(connStr string) error {

	var err error
	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Error connecting to the database.", err)
	}
	if err = DB.Ping(); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

	// err = DB.Ping()

	// if err != nil {
	// 	log.Fatal("Database Unreachable")
	// }
	log.Println("Database connected successfully")
	return nil
}

