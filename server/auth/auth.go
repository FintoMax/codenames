package auth

import (
	"codenames/server/database"
	"codenames/server/domain"
	"codenames/server/models"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID int64
	jwt.RegisteredClaims
}
type AuthService struct {
	db *sqlx.DB
}

func generateToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
func (s *AuthService) CreateUser(username string, password string) (string, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return "", err
	}
	_, err = database.GetUserByUsername(s.db, username)
	if err == nil {
		return "", domain.ErrUserAlreadyExists
	}
	id, err := database.CreateUser(s.db, username, passwordHash)
	if err != nil {
		return "", err
	}
	token, err := generateToken(int64(id))
	if err != nil {
		return "", err
	}
	return token, nil
}
func (s *AuthService) Login(username string, password string) (string, error) {
	user, err := database.GetUserByUsername(s.db, username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidPassword
	}
	token, err := generateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil

}
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
