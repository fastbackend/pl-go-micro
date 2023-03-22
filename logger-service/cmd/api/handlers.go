package main

import (
	"net/http"

	"github.com/fastbackend/pl-go-micro/logger-service/data"
)

// Funkcja zapisuje zdarzenia do bazy mongoDB
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {

	// Definicja struktury payloadu
	var requestPayload struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	// Próba odczytu danych uwierzytelniających z ciała żądania
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		// Jeżeli nie udana, to zwracamy błąd
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Próba dodania danych o zdarzeniu
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		// Jeżeli nie udana, to zwracamy błąd
		app.errorJSON(w, err)
		return
	}

	// Utworzenie obiektu jsonResponse z informacjami o rejestracji
	response := jsonResponse{
		Error:   false,
		Message: "Zarejestrowano zdarzenie",
	}

	// Zwrócenie danych JSON klientowi
	app.writeJSON(w, http.StatusAccepted, response)

}
