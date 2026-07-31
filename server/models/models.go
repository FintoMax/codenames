package models

import "sync"

type User struct {
	ID           int64  `db:"id"`
	Name         string `db:"name"`
	PasswordHash string `db:"password_hash"`
}

type Users struct {
	mu    sync.RWMutex
	users []User
}

type Lobby struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	HostID       int64  `json:"host_id" db:"host_id"`
	CreatedAt    string `json:"created_at" db:"created_at"`
}
