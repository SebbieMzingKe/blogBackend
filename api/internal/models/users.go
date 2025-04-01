package models

import "database/sql"

type User struct {
	ID       int            `json:"id"`
	Email    string         `json:"email"`
	Password string         `json:"-"`
	Name     sql.NullString `json:"name"`
}
