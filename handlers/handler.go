package handlers

import (
	"context"
	"example/GoApiTemplate/models"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// Handlers should invoke functions to get the data without accessing the database

func AddItem(context echo.Context) error {
	item := new(models.Item)

	err := context.Bind(&item)

	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid data.")
	}

	newItem := models.NewItem(item.Name, item.Price)

	models.Items = append(models.Items, newItem)

	return context.JSON(http.StatusCreated, newItem.ID)
}

func DeleteItem(context echo.Context) error {
	id, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid id.")
	}

	for index, i := range models.Items {
		if i.ID == id {
			models.Items = append(models.Items[:index], models.Items[index+1:]...)

			return context.JSON(http.StatusNoContent, i.ID)
		}
	}

	return context.JSON(http.StatusNotFound, "Item not found.")
}

func EditItem(context echo.Context) error {
	id, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid id.")
	}

	item := new(models.Item)

	err = context.Bind(&item)

	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid data.")
	}

	for index := range models.Items {
		if models.Items[index].ID == id {
			models.Items[index].Name = item.Name
			models.Items[index].Price = item.Price

			return context.JSON(http.StatusOK, models.Items[index])
		}
	}

	return context.JSON(http.StatusNotFound, "Item not found.")
}

func GetItem(context echo.Context) error {
	idStr, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid request")
	}
	for _, item := range models.Items {
		if item.ID == idStr {
			return context.JSON(http.StatusOK, item)
		}
	}

	return context.JSON(http.StatusNotFound, "Not found")
}

type ItemHandler struct {
	DB *pgxpool.Pool
}

func (_itemHandler *ItemHandler) GetItems(c echo.Context) error {
	rows, err := _itemHandler.DB.Query(context.Background(), "SELECT * FROM Item")

	if err != nil {
		return c.JSON(http.StatusBadGateway, "Bad gateway?")
	}

	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var item models.Item

		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			log.Printf("Error scanning row: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error."})
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, items)
}
