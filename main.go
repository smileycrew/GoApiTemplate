package main

import (
	"context"
	"example/GoApiTemplate/handlers"
	"example/GoApiTemplate/models"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	models.Items = append(models.Items, models.NewItem("Item 1", 4))
	models.Items = append(models.Items, models.NewItem("Item 2", 6))

	app := echo.New()
	app.Use(middleware.Logger())

	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading env file.")
	}

	// connection string needs to be in os.Getenv
	connStr := os.Getenv("DB_CONNECTION_STRING")

	// connect to the database
	dbPool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v/n", err)
		os.Exit(1)
	}

	defer dbPool.Close()

	// verify the connection

	conn := &handlers.ItemHandler{DB: dbPool}

	app.DELETE("/items/:id", handlers.DeleteItem)
	app.GET("/items", conn.GetItems)
	app.GET("/items/:id", handlers.GetItem)
	app.POST("/items", handlers.AddItem)
	app.PUT("/items/:id", handlers.EditItem)

	app.Logger.Fatal(app.Start(":8080"))
}
