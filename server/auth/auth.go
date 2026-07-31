package auth

import (
	"codenames/server/database"
	"codenames/server/domain"
	"errors"
	"os"
	"sync"
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
	db     *sqlx.DB
	tokens map[string]int64
	mu     sync.Mutex
}

func NewAuthService(db *sqlx.DB) *AuthService {
	service := &AuthService{
		db:     db,
		tokens: make(map[string]int64),
	}
	go service.cleanuploop()
	return service
}

func (s *AuthService) cleanuploop() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		s.cleanupExpiredTokens()
	}
}
func (s *AuthService) cleanupExpiredTokens() {
	s.mu.Lock()
	defer s.mu.Unlock()
	timeNow := time.Now()
	for token := range s.tokens {
		parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			delete(s.tokens, token)
			continue
		}
		if claims, ok := parsed.Claims.(*Claims); ok {
			if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(timeNow) {
				delete(s.tokens, token)
			}
		}
	}
}

func (s *AuthService) generateToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	s.tokens[token] = userID
	s.mu.Unlock()
	return token, err
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
	token, err := s.generateToken(int64(id))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) Login(username string, password string) (string, error) {
	user, err := database.GetUserByUsername(s.db, username)
	if err != nil {
		return "", domain.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidPassword
	}
	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Logout(tokenString string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tokens, tokenString)
	return nil
}

func (s *AuthService) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	s.mu.Lock()
	_, exists := s.tokens[tokenStr]
	s.mu.Unlock()
	if !exists {
		return nil, errors.New("token not found")
	}
	return claims, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *AuthService) GetDB() *sqlx.DB {
	return s.db
}
