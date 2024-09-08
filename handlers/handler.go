package handlers

import (
	"example/GoApiTemplate/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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

func GetItems(context echo.Context) error {
	items := models.Items

	return context.JSON(http.StatusOK, items)
}

func ItemHandler(writer http.ResponseWriter, request *http.Request) {
	// depending on the request method it will route to the correct function
	switch request.Method {
	case http.MethodGet:
		// GetItem(request, writer)
	default:
		http.Error(writer, "Method not allowed.", http.StatusMethodNotAllowed)
	}
}
