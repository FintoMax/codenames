package database

import (
	"codenames/server/config"
	_ "embed"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

//go:embed init.sql
var initSql string

func NewDB(cfg *config.Config) (*sqlx.DB, error) {
	connStr, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	db, err := sqlx.Open("pgx", connStr.ConnString)
	if err != nil {
		log.Fatal("Error connecting to database: " + err.Error())
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}
	return db, nil
}
