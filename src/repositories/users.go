package repositories

import (
	"api/src/models"
	"database/sql"
)

// Represent a user repository
type Users struct {
	db *sql.DB
}

// Create a user repository
func NewUsersRepository(db *sql.DB) *Users {
	return &Users{db}
}

// Inserts a user into the database
func (repository Users) Create(user models.User) (uint64, error) {
	statement, err := repository.db.Prepare(
		"INSERT INTO users (name, nick, email, password) values (?, ?, ?, ?)",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(user.Name, user.Nick, user.Email, user.Password)
	if err != nil {
		return 0, err
	}
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(lastInsertedID), nil
}
