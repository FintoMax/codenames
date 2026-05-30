package models

import "sync"

type User struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	PasswordHash string `db:"password_hash"`
}

type Users struct {
	mu    sync.RWMutex
	users []User
}
