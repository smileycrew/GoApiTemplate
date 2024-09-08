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

func (itemHandler *ItemHandler) Controller(cntx echo.Context) error {
	switch cntx.Request().Method {
	case http.MethodDelete:
		itemHandler.deleteItem(cntx)
	case http.MethodGet:
		itemHandler.getItems(cntx)
	case http.MethodPost:
		return cntx.JSON(http.StatusOK, "POST")
	case http.MethodPut:
		return cntx.JSON(http.StatusOK, "PUT")
		// what error should go here???? defaultl i think
	}

	return cntx.JSON(http.StatusOK, "NO ROUTE FOUND")
}

func (itemHandler *ItemHandler) addItem(context echo.Context) error {
	item := new(models.Item)

	err := context.Bind(&item)

	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid data.")
	}

	newItem := models.NewItem(item.Name, item.Price)

	models.Items = append(models.Items, newItem)

	return context.JSON(http.StatusCreated, newItem.ID)
}

func (itemHandler *ItemHandler) deleteItem(cntx echo.Context) error {
	id, err := strconv.Atoi(cntx.Param("id"))

	log.Print(cntx.Param("id"))

	if err != nil {
		return cntx.JSON(http.StatusBadRequest, "Invalid id.")
	}

	rows, err := itemHandler.DB.Query(context.Background(), "DELETE FROM Item WHERE id = ?", id)

	if err != nil {
		return cntx.JSON(http.StatusNotFound, "Item not found.")
	}

	return cntx.JSON(http.StatusOK, rows)
}

func (itemHandler *ItemHandler) editItem(context echo.Context) error {
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

func (itemHandler *ItemHandler) getItems(c echo.Context) error {
	// i need a way to get the DB from ItemHandlerStrut here

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
