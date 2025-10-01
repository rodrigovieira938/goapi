package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rodrigovieira938/goapi/config"
)

func Connect(dbConfig config.DatabaseConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.Hostname, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Database)
	return sql.Open("postgres", psqlInfo)
}
func GetUserById(db *sql.DB, id int) *sql.Row {
	row := db.QueryRow("SELECT * from \"user\" WHERE id = $1", id)
	if row == nil {
		return nil
	}
	return row
}
func GetUserByEmail(db *sql.DB, email string) *sql.Row {
	row := db.QueryRow("SELECT * from \"user\" WHERE email = $1", email)
	return row
}
func GetCarById(db *sql.DB, id int) *sql.Row {
	row := db.QueryRow("SELECT * from car WHERE id = $1", id)
	if row == nil {
		return nil
	}
	return row
}
