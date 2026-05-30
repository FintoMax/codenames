package auth

import (
	"codenames/config"
	"database/sql"
)

type AuthService struct {
	cfg *config.Config
	db  *sql.DB
}

func NewAuthService(cfg *config.Config, db *sql.DB) *AuthService {
	return &AuthService{cfg: cfg, db: db}
}
