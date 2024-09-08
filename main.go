package main

import (
	"context"
	"example/GoApiTemplate/handlers"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

func main() {
	// load env files
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env files: %v", err)
	}

	// initialize echo
	app := echo.New()
	app.Use(middleware.Logger())

	fmt.Println("Environment variables loaded successfully.")

	connStr := os.Getenv("DB_CONNECTION_STRING")

	// creating a new connection pool to the database
	dbPool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)

		os.Exit(1)
	}

	defer dbPool.Close()

	// verify the connection
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	h := handlers.NewItemHandler(dbPool)

	// may need to add the standard get in here
	app.Any("/items*", h.Controller)

	app.Logger.Fatal(app.Start(":8080"))
}
