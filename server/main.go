package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}
	defer db.Close()

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, "Hello, World!")
	})
	if err := e.Start(":8080"); err != nil {
		e.Logger.Info("shutting down the server")
	}

	e.POST("/auth/signup", createUser)
	e.POST("/auth/login", login)
	e.POST("/lobby/create", createLobby)
	e.POST("/lobby/join", joinLobby)
}
