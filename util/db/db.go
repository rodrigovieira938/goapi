package db

import (
	"database/sql"
	"fmt"

	"github.com/rodrigovieira938/goapi/config"

	_ "github.com/lib/pq"
)

func Connect(dbConfig config.DatabaseConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.Hostname, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Database)
	return sql.Open("postgres", psqlInfo)
}
