package main

import (
	"net/http"
)

// Funkcja wysyła mail
func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {

	// Definicja struktury payloadu
	var requestPayload struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	// Próba odczytu danych uwierzytelniających z ciała żądania
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		// Jeżeli nie udana, to zwracamy błąd
		app.errorJSON(w, err)
		return
	}

	// Próba wysyłania maila
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}
	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		// Jeżeli nie udana, to zwracamy błąd
		app.errorJSON(w, err)
		return
	}

	// Utworzenie obiektu jsonResponse z informacjami o rejestracji
	response := jsonResponse{
		Error:   false,
		Message: "Wiadomość wysłana do " + requestPayload.To,
	}

	// Zwrócenie danych JSON klientowi
	app.writeJSON(w, http.StatusAccepted, response)

}
