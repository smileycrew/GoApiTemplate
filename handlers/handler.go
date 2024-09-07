package handlers

import (
	"encoding/json"
	"example/GoApiTemplate/models"
	"net/http"
)

var (
	headerKey   = "Content-Type"
	headerValue = "application/json"
)

var Items []models.Item

func GetItem(request *http.Request, writer http.ResponseWriter) {
	// sets the http response to be returned back to the client
	writer.Header().Set(headerKey, headerValue)

	// ecodes the items into JSON and writes the data directly into the writer
	json.NewEncoder(writer).Encode(Items)
}

func GetItems(writer http.ResponseWriter, request *http.Request) {
	// sets the https response to be returned back to the client
	writer.Header().Set(headerKey, headerValue)

	// encodes the items into JSON and writes the data directly into the writer
	json.NewEncoder(writer).Encode(Items)
}

func ItemHandler(writer http.ResponseWriter, request *http.Request) {
	// depending on the request method it will route to the correct function
	switch request.Method {
	case http.MethodGet:
		GetItem(request, writer)
	default:
		http.Error(writer, "Method not allowed.", http.StatusMethodNotAllowed)
	}
}
