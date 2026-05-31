package auth

import (
	"codenames/server/database"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secret-key")

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("Invalid token")
	}
	return nil
}

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
