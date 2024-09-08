package main

import (
	"example/GoApiTemplate/handlers"
	"example/GoApiTemplate/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	models.Items = append(models.Items, models.NewItem("Item 1", 4))
	models.Items = append(models.Items, models.NewItem("Item 2", 6))

	app := echo.New()
	app.Use(middleware.Logger())

	app.DELETE("/items/:id", handlers.DeleteItem)
	app.GET("/items", handlers.GetItems)
	app.GET("/items/:id", handlers.GetItem)
	app.POST("/items", handlers.AddItem)
	app.PUT("/items/:id", handlers.EditItem)

	app.Logger.Fatal(app.Start(":8080"))
}
