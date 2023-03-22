package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Stała określająca numer portu
const webPort = "80"

// Definicja struktury Config
type Config struct {
	Mailer Mail
}

func main() {

	// Logowanie informacji o starcie usługi Mailer'a na porcie 80
	log.Printf("Uruchamianie usługi wysyłania maili na porcie %s\n", webPort)

	// Inicjalizacja instancji aplikacji z konfiguracją Config
	app := Config{
		Mailer: createMail(),
	}

	// Definicja serwera HTTP
	srv := &http.Server{
		// Adres, na którym serwer będzie nasłuchiwać żądania HTTP
		Addr: fmt.Sprintf(":%s", webPort),
		// Funkcja obsługi żądań przychodzących na serwer
		Handler: app.routes(),
	}

	// Nasłuchiwanie na adresie i obsługa żądań HTTP
	err := srv.ListenAndServe()
	if err != nil {
		// W przypadku błędu zakończ program i zwróć błąd
		log.Panic(err)
	}
}

// Konfiguracja serwera mailowego
func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	mail := Mail{
		// zweryfikować MAIL_DOMAIN, w kontekście MAIL_DRIVER
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	// Zwraca konfigurację obiektu Mail
	return mail
}

// powyższe parametry w kontekście docker-compose.yml (z project)
