package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Definicja struktury jsonResponse,
// która jest używana do tworzenia odpowiedzi JSON
type jsonResponse struct {
	// Pole "Error" określa, czy wystąpił błąd podczas przetwarzania żądania
	Error bool `json:"error"`
	// Pole "Message" zawiera wiadomość zwrotną dla żądania
	Message string `json:"message"`
	// Pole "Data" zawiera opcjonalne dane zwrotne z żądania
	Data any `json:"data,omitempty"`
}

// readJSON próbuje odczytać r.Body żądania i konwertuje je na JSON
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	// Ustawienie maksymalnej liczby bajtów, jakie mogą być przeczytane z r.Body
	maxBytes := 1048576 // zabezpieczenie: jeden megabajt
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Dekodowanie r.Body do wartości data
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// Sprawdzenie, czy r.Body zawiera tylko jedną wartość JSON
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body może zawierać tylko jedną wartość JSON")
	}

	// Jeżeli wszystko jest w porządku, to zwracamy nil
	return nil
}

// writeJSON przyjmuje kod statusu odpowiedzi oraz dowolne dane i zapisuje odpowiedź json dla klienta
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	// Przekształć obiekt jsonResponse na JSON z wcięciami i przypisz do zmiennej out
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Jeśli w "headers" jest przynajmniej jeden nagłówek,
	// przepisz jego wartości do odpowiedzi HTTP
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	// Ustaw nagłówek Content-Type na application/json
	w.Header().Set("Content-Type", "application/json")
	// Ustaw kod statusu odpowiedzi
	w.WriteHeader(status)
	// Zapisz odpowiedź JSON do ResponseWriter
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON przyjmuje błąd i opcjonalnie kod statusu odpowiedzi,
// generuje i wysyła odpowiedź na błąd w formacie json
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {

	// Domyślnie kod statusu to 400 (Bad Request)
	statusCode := http.StatusBadRequest

	// Jeśli podano kod statusu, przypisz go do zmiennej "statusCode"
	if len(status) > 0 {
		statusCode = status[0]
	}

	// Utwórz strukturę "jsonResponse" zawierającą informacje o błędzie
	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	// Wywołaj funkcję "writeJSON()", która wygeneruje i wyśle odpowiedź na błąd w formacie JSON
	return app.writeJSON(w, statusCode, payload)
}
