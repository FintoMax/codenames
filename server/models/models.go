package models

type User struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	PasswordHash string `db:"password_hash"`
}
