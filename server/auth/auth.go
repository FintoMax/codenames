package auth

import (
	"codenames/server/database"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func IfUserExists(db *sqlx.DB, username string) bool {
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE name = $1)", username)
	if err != nil {
		return false
	}
	return exists
}

func AuthenticateUser(db *sqlx.DB, username, password string) error {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return err
	}
	if IfUserExists(db, username) {
		fmt.Printf("User %s already exists\n", username)
	}
	_, err = database.CreateUser(db, username, passwordHash)
	if err != nil {
		return err
	}
	return nil
}

func VerifyUser(db *sqlx.DB, username, password string) error {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return err
	}
	if !IfUserExists(db, username) {
		fmt.Printf("User %s does not exist\n", username)
	}
	user, err := database.GetUserByUsername(db, username)
	if err != nil {
		return err
	}
	if user.PasswordHash != passwordHash {
		fmt.Printf("Password is incorrect\n", username)
	}
	return nil
}
