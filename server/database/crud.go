package database

import (
	"codenames/server/models"
	"errors"

	"github.com/jmoiron/sqlx"
)

func CreateUser(db *sqlx.DB, name, passwordHash string) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO users (name, password_hash) VALUES ($1, $2) RETURNING id", name, passwordHash).Scan(&id)
	return id, err
}

func GetUserbyID(db *sqlx.DB, id int) (*models.User, error) {
	user := &models.User{}
	err := db.Get(user, "SELECT id, name, password_hash FROM users WHERE id = $1", id)
	return user, err
}

func GetUserByUsername(db *sqlx.DB, username string) (*models.User, error) {
	user := &models.User{}
	err := db.Get(user, "SELECT id, name, password_hash FROM users WHERE username = $1", username)
	return user, err
}

func GetAllUsers(db *sqlx.DB) ([]models.User, error) {
	users := []models.User{}
	err := db.Select(&users, "SELECT id, name, password_hash FROM users")
	return users, err
}

func UpdateUser(db *sqlx.DB, user *models.User) error {
	res, err := db.Exec("UPDATE users SET name = $1 WHERE id = $2", user.Name, user.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("User does not exist")
	}
	return nil
}

func DeleteUser(db *sqlx.DB, id int) error {
	res, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("User does not exist")
	}
	return nil
}
