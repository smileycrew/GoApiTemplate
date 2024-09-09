package handlers

import (
	"context"
	"example/GoApiTemplate/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

// Handlers should invoke functions to get the data without accessing the database

type ItemHandler struct {
	DB *pgxpool.Pool
}

func NewItemHandler(db *pgxpool.Pool) *ItemHandler {
	return &ItemHandler{DB: db}
}

func (itemHandler *ItemHandler) AddItem(cntx echo.Context) error {
	item := new(models.Item)

	err := cntx.Bind(&item)

	if err != nil {
		return cntx.JSON(http.StatusBadRequest, "Invalid data.")
	}

	newItem := models.NewItem(item.Name, item.Price)

	if err := itemHandler.DB.QueryRow(context.Background(), "INSERT INTO Item(name, price) VALUES($1, $2) RETURNING id, name, price", item.Name, item.Price).Scan(&newItem.ID, &newItem.Name, &newItem.Price); err != nil {
		return cntx.JSON(http.StatusInternalServerError, "Internal server error.")
	}

	return cntx.JSON(http.StatusCreated, newItem)
}

func (itemHandler *ItemHandler) DeleteItem(cntx echo.Context) error {
	id, err := strconv.Atoi(cntx.Param("id"))

	log.Printf("The param id is: %v", id)

	if err != nil {
		return cntx.JSON(http.StatusBadRequest, "Invalid id.")
	}

	rows, err := itemHandler.DB.Exec(context.Background(), "DELETE FROM Item WHERE id = $1", id)

	if err != nil {
		return cntx.JSON(http.StatusNotFound, "Item not found.")
	}

	return cntx.JSON(http.StatusOK, rows)
}

func (itemHandler *ItemHandler) EditItem(cntx echo.Context) error {
	id, err := strconv.Atoi(cntx.Param("id"))

	if err != nil {
		return cntx.JSON(http.StatusBadRequest, "Invalid id.")
	}

	item := new(models.Item)

	if err := cntx.Bind(&item); err != nil {
		return cntx.JSON(http.StatusBadRequest, "Invalid data.")
	}

	result, err := itemHandler.DB.Exec(context.Background(), "UPDATE Item SET name = $1, price = $2 WHERE id = $3", item.Name, item.Price, id)

	if err != nil {
		return cntx.JSON(http.StatusNotFound, "Item not found.")
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return cntx.JSON(http.StatusInternalServerError, "Could not update item.")
	}

	rows, err := itemHandler.DB.Query(context.Background(), "SELECT * FROM Item WHERE id = $1", id)

	if err != nil {
		return cntx.JSON(http.StatusNotFound, "Item not found after update.")
	}

	for rows.Next() {
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			log.Printf("Error scanning row: %v", err)

			return cntx.JSON(http.StatusInternalServerError, "Internal server errors.")
		}
	}

	return cntx.JSON(http.StatusOK, item)
}

func (itemHandler *ItemHandler) GetItem(cntx echo.Context) error {
	id, err := strconv.Atoi(cntx.Param("id"))

	if err != nil {
		return cntx.JSON(http.StatusBadRequest, "Invalid request")
	}

	rows, err := itemHandler.DB.Query(context.Background(), "SELECT * FROM Item WHERE id = $1", id)

	if err != nil {
		return cntx.JSON(http.StatusNotFound, "Not found")
	}

	defer rows.Close()

	item := new(models.Item)

	for rows.Next() {
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			log.Printf("Error scanning row: %v:", err)

			return cntx.JSON(http.StatusInternalServerError, "Internal server error.")
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)

		return cntx.JSON(http.StatusInternalServerError, "Internal server error.")
	}

	return cntx.JSON(http.StatusOK, item)
}

func (itemHandler *ItemHandler) GetItems(c echo.Context) error {
	rows, err := itemHandler.DB.Query(context.Background(), "SELECT * FROM Item")

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
