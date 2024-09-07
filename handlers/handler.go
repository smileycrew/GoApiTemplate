package handlers

import (
	"example/GoApiTemplate/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// Handlers should invoke functions to get the data without accessing the database

func GetItem(context echo.Context) error {
	idStr, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid request")
	}
	for _, item := range models.Items {
		if item.Id == idStr {
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
