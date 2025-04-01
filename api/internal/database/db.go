package database

import (
	"blogBackend/internal/models"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func CreateUserTable()  {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		name TEXT
	);`

	if _, err := DB.Exec(query); err != nil {
		log.Fatal("Error creating user tables", err)
	}
}

func CreateUser(email, hashedPassword, name string) error {
	_, err := DB.Exec("INSERT INTO users (email, password, name) VALUES ($1, $2, $3)", email, hashedPassword, name)
	return err
}
	
func GetUserByEmail(email string) (models.User, error) {
	var u models.User
	err := DB.QueryRow("SELECT id, email, password, name FROM users WHERE email = $1", email).Scan(&u.ID, &u.Email, &u.Password, &u.Name)
	
	return u, err
}

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
	log.Println("Database connected successfully")
	return nil
}

