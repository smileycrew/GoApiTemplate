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

	connStr := os.Getenv("DB_CONNECTION_STRING")

	fmt.Println("Environment variables loaded successfully.")

	// initialize echo
	app := echo.New()
	app.Use(middleware.Logger())

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

	itemHandler := handlers.NewItemHandler(dbPool)

	app.DELETE("/items/:id", itemHandler.DeleteItem)
	app.GET("/items", itemHandler.GetItems)
	app.GET("/items/:id", itemHandler.GetItem)
	app.POST("/items", itemHandler.AddItem)
	app.PUT("/items/:id", itemHandler.EditItem)

	app.Logger.Fatal(app.Start(":8080"))
}
