package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fastbackend/pl-go-micro/broker-service/event"
)

// Definicja struktury Request
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

// Definicja struktury AuthPayload dla RequestPayload
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Definicja struktury LogPayload dla RequestPayload
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Definicja struktury MailPayload dla RequestPayload
type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

/*
Testowanie działania Broker'a
*/

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	// Utwórz obiekt jsonResponse
	payload := jsonResponse{
		Error:   false,
		Message: "Broker potwierdza gotowość do przyjmowania żądań.",
	}

	// Zwracamy wiadomość (payload) z kodem 200
	_ = app.writeJSON(w, http.StatusOK, payload)
}

/*
Funkcja obsługi żądań dla identyfikacji mikroserwisu i dalszego wywołania
*/

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	// Deklaracja zmiennej dla payloadu żądania
	var requestPayload RequestPayload

	// Wywołanie metody readJSON
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		// W przypadku błędu odczytu zwórcenie błędu i przerwanie
		app.errorJSON(w, err)
		return
	}

	// Instrukcja switch dla rozdzielenia żądań, zależnie od wywołania
	switch requestPayload.Action {
	// Dla "auth" wywołaj poniższą funkcję authenticate
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	// Dla "log" wywołaj poniższą funkcję logItem
	case "log":
		app.logEventViaRabbit(w, requestPayload.Log) //RMQ
		//app.logItem(w, requestPayload.Log)
	// Dla "mail" wywołaj poniższą funkcję sendMail
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	// Jeżeli typ wywołania nie został wywołany zwróć błąd z komunikatem jak poniżej
	default:
		app.errorJSON(w, errors.New("niezidentyfikowane żądanie"))
	}

}

/*
Funkcja wywołująca mikroserwis autentykacji auth-service
*/

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {

	// Tworzenie JSON, który wyślemy do mikroserwisu autentykacji
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// Wywołanie serwisu auth-service
	authServiceURL := "http://auth-service/authenticate"
	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Utworzenie klienta HTTP i wysłanie żądania (request) do serwisu autentykacji
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		// W przypadku błędu:
		app.errorJSON(w, err)
		return
	}
	// Zamknięcie response.Body po zakończeniu działania funkcji,
	// aby zwolnić zasoby sieciowe:
	defer response.Body.Close()

	// Upewniamy się, że otrzymamy poprawny kod statusu
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("nieprawidłowe dane uwierzytelniające"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("błąd wywołania usługi auth (zły mail, brak tablicy etc.)"))
		return
	}

	// Do tej zmiennej będziemy przypisywać response.Body
	var jsonFromService jsonResponse

	// Dekodowanie JSON z auth-service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Jeżeli wystąpił błąd autentyfikacji
	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// Jeżeli wszystko poszło poprawnie, to:
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Weryfikacja z Postgres. Dane logowania prawidłowe. "
	payload.Data = jsonFromService.Data

	// i zwracamy wynik w postaci JSON'a
	app.writeJSON(w, http.StatusAccepted, payload)

}

/*
Funkcja wywołująca mikroserwis logger-service
*/

// Nie używana w związku z testowaniem RabbitMQ
func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {

	// Tworzenie JSON, który wyślemy do mikroserwisu dziennika zdarzeń
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// Wywołanie serwisu logger-service
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Utworzenie klienta HTTP i wysłanie żądania (request) do serwisu dziennika zdarzeń
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		// W przypadku błędu:
		app.errorJSON(w, err)
		return
	}
	// Zamknięcie response.Body po zakończeniu działania funkcji,
	// aby zwolnić zasoby sieciowe:
	defer response.Body.Close()

	// Jeżeli nie otrzymaliśmy poprawnego kodu statusu
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	// Jeżeli wszystko poszło poprawnie, to:
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Testowa rejestracja wykonana. Dopisano do mongoDB."

	// i zwracamy wynik w postaci JSON'a
	app.writeJSON(w, http.StatusAccepted, payload)

}

/*
Funkcja wywołująca mikroserwis mailer-service
*/

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {

	// Tworzenie JSON, który wyślemy do mikroserwisu wysyłania maili
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// Wywołanie serwisu mailer-service
	mailServiceURL := "http://mailer-service/send"
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Utworzenie klienta HTTP i wysłanie żądania (request) do serwisu wysyłania maili
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		// W przypadku błędu:
		app.errorJSON(w, err)
		return
	}
	// Zamknięcie response.Body po zakończeniu działania funkcji,
	// aby zwolnić zasoby sieciowe:
	defer response.Body.Close()

	// Jeżeli nie otrzymaliśmy poprawnego kodu statusu
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("błąd wywołania usługi mailowej (podaj prawidłowe parametry smtp)"))
		return
	}

	// Jeżeli wszystko poszło poprawnie, to:
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Wyślij wiadomość do " + msg.To

	// i zwracamy wynik w postaci JSON'a
	app.writeJSON(w, http.StatusAccepted, payload)

}

/*
Funkcja przykładowe dla obsługi RabbitMQ
*/

// RMQ: logEventViaRabbit rejestruje zdarzenie za pomocą usługi logger-service.
// Wykonuje on połączenie, wypychając dane do RabbitMQ
func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Rejestracja zdarzenia do mongoDB poprzez RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)
}

// RMQ: pushToQueue przesyła wiadomość do RabbitMQ
func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
