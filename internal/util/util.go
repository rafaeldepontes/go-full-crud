package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const BrazilianDateTimeFormat = "02/01/2006 15:04:05"

var ErrorBlankId = errors.New("Id is required")
var ErrorUserNotFound = errors.New("User not found")
var ErrorBlankUsername = errors.New("Username is required")
var ErrorMethodNotAllowed = errors.New("Method not allowed")
var ErrorUserCredentials = errors.New("Invalid username or password")

type Error struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

var (
	BadRequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	RequestErrorHandler = func(w http.ResponseWriter, err error, status int) {
		writeError(w, err.Error(), status)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An unexpected Error Occurred.", http.StatusInternalServerError)
	}
)

func writeError(w http.ResponseWriter, message string, status int) {
	var timestamp string = time.Now().Format(BrazilianDateTimeFormat)
	resp := Error{
		Status:    status,
		Message:   message,
		Timestamp: timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error(err)
	}
}
