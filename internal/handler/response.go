package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	log.Panic(message)
	w.WriteHeader(statusCode)
	result, _ := json.Marshal(errorResponse{message})
	_, _ = w.Write(result)
}
