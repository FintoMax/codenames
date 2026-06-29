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
