package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

/*

Funkcja Authenticate jest używana do uwierzytelniania użytkownika
na podstawie otrzymanych danych uwierzytelniających, tj. adresu email i hasła.

1.
Pierwszym krokiem jest próba odczytu tych danych z ciała żądania.
Jeśli wystąpi błąd podczas próby odczytu, funkcja errorJSON zostanie wywołana, aby zwrócić błąd klientowi.

2.
Następnie funkcja GetByEmail() jest używana do pobrania użytkownika z bazy danych na podstawie podanego adresu email.
Jeśli nie uda się pobrać użytkownika, errorJSON zostanie wywołana z błędem "nieprawidłowe dane uwierzytelniające".

3.
Następnie sprawdzane jest, czy podane hasło odpowiada hasłu użytkownika.
Jeśli hasło jest nieprawidłowe lub wystąpił błąd podczas sprawdzania,
errorJSON zostanie wywołana z błędem "nieprawidłowe dane uwierzytelniające".

4.
Jeśli użytkownik jest poprawnie uwierzytelniony,
zostanie utworzony obiekt jsonResponse z informacjami o zalogowanym użytkowniku.

5.
Ostatecznie, dane JSON zostaną zwrócone klientowi za pomocą funkcji writeJSON().

*/

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	// Definicja struktury payloadu
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Próba odczytu danych uwierzytelniających z ciała żądania
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		// Jeżeli nie udana, to zwracamy błąd
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Pobranie użytkownika z bazy danych na podstawie podanego adresu email
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		// Jeżeli nie udane, to zwracamy błąd
		app.errorJSON(w, errors.New("nieprawidłowe dane uwierzytelniające (email)"), http.StatusBadRequest)
		return
	}

	// Sprawdzenie, czy podane hasło odpowiada hasłu użytkownika (pasuje do hasha'a)
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		// Jeżeli nie pasuję, to zwracamy błąd
		app.errorJSON(w, errors.New("nieprawidłowe dane uwierzytelniające (hasło)"), http.StatusBadRequest)
		return
	}

	// Ta sekcja została dodana dla współpracy z logger-service
	err = app.logRequest("autentykacja", fmt.Sprintf("%s zautoryzowany", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Utworzenie obiektu jsonResponse z informacjami o zalogowanym użytkowniku
	response := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Zalogowany użytkownik %s", user.Email),
		Data:    user,
	}

	// Zwrócenie danych JSON klientowi
	app.writeJSON(w, http.StatusAccepted, response)

}

// Ta funkcja została dodana w celu rejestrowania w logger-service żądań autentykacji
func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
