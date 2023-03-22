package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Funkcja routes zwraca router HTTP jako typ http.Handler,
// który jest gotowy do obsługi żądań HTTP dla Logger'a
func (app *Config) routes() http.Handler {
	// Utwórz nowy router HTTP
	mux := chi.NewRouter()

	// Zabezpieczenia CORS określają,
	// kto może łączyć się z serwerem i jakie zapytania są dozwolone
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, //ciasteczka, tokeny lub certyfikaty SSL
		MaxAge:           300,  //buforowanie CORS na określony czas (300 sek.)
	}))

	// Użyj funkcji middleware.Heartbeat(),
	// która dodaje funkcjonalność śledzenia stanu serwera
	mux.Use(middleware.Heartbeat("/ping"))

	// To jest wywoływane przez Brokera w jego funkcji HandleSubmission
	mux.Post("/log", app.WriteLog)

	// Zwróć router HTTP jako typ http.Handler,
	// który jest gotowy do obsługi żądań HTTP
	return mux
}
