package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Rejestracja funkcji obsługi dla ścieżki podstawowej ("/")
	// Funkcja obsługi wywołuje funkcję render() z argumentami:
	// http.ResponseWriter i nazwą szablonu GoTemplate "test.page.gohtml"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "broker.page.gohtml")
	})

	// Wyświetlenie komunikatu w CLI o uruchamianiu serwera na porcie 80
	fmt.Println("Uruchamianie usługi front-endu na porcie 80")

	// Uruchomienie serwera web na porcie 80, przy użyciu funkcji http.ListenAndServe()
	// W przypadku wystąpienia błędu, funkcja log.Panic() zostanie wywołana z argumentem err,
	// co spowoduje zakończenie działania programu.
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

// Funkcja obsługi szablonów GoTemplates, w tym renderowanie
func render(w http.ResponseWriter, t string) {
	// Określenie ścieżek do plików szablonów *.gohtml
	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	// Utworzenie slice'a zawierającego ścieżkę do szablonu, a także ścieżki do plików szablonów częściowych
	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))

	// Dodanie wszystkich elementów ze slice'u `partials` do slice'u `templateSlice`
	templateSlice = append(templateSlice, partials...)

	// Parsowanie plików szablonów i renderowanie wyniku na obiekcie http.ResponseWriter
	// Tworzenie obiektu *template i przechowywanie w pamięci operacyjnej dla szybszego wykonania
	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		// W przypadku wystąpienia błędu podczas parsowania szablonów, wyświetlenie błędu w przeglądarce klienta i zakończenie funkcji.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Korzystanie z obiektu tmpl w celu wyświetlenia strony w przeglądarce klienta
	if err := tmpl.Execute(w, nil); err != nil {
		// W przypadku wystąpienia błędu podczas renderowania szablonów, wyświetlenie błędu w przeglądarce klienta
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
