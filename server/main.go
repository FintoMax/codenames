package main

import (
	"codenames/server/auth"
	"codenames/server/handler"
	"codenames/server/middleware"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("CONN_STRING")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}
	defer db.Close()

	authService := auth.NewAuthService(db)
	e := echo.New()

	e.Use(echoMiddleware.RequestLogger())
	e.Use(echoMiddleware.Recover())

	e.POST("/auth/signup", handler.CreateUser(authService))
	e.POST("/auth/login", handler.Login(authService))

	e.Use(middleware.AuthMiddleware(authService))

	e.GET("/", func(c *echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/auth/login")
	})
	e.POST("/lobby/create", handler.CreateLobby(authService))
	e.POST("/lobby/join", handler.JoinLobby(authService))

	if err := e.Start(":8080"); err != nil {
		e.Logger.Info("shutting down the server")
	}

}
