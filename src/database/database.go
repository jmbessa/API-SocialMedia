package database

import (
	"api/src/config"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Open the connection with the database
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.ConnectionString)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)

	db.SetMaxIdleConns(5)

	db.SetConnMaxIdleTime(15 * time.Minute)

	db.SetConnMaxLifetime(2 * time.Hour)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
