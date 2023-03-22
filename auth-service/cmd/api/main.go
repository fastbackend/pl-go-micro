package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fastbackend/pl-go-micro/auth-service/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Stała określająca numer portu
const webPort = "80"

// Licznik dla odliczania ilości prób połączeń z DB
var counts int64

// Definicja struktury Config
type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	// Logowanie informacji o starcie usługi autentykacji na porcie 80
	log.Printf("Uruchamianie usługi autentykacji na porcie %s\n", webPort)

	// Tworzenie połączenia do bazy danych
	conn := connectToDB()
	if conn == nil {
		log.Panic("Nie można połączyć się z Postgres!")
	}

	// Inicjalizacja instancji aplikacji z konfiguracją Config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
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

// Funkcja do otwierania połączenia z bazą z wykorzystaniem
// sterownika "pgx" i parametrów DSN podawnaych z connectToDB
// (ta funkcja jest pomocnicza dla connectToDB)
func openDB(dsn string) (*sql.DB, error) {
	// Wykorzystanie sterownika pgx i parametrów dsn
	// dla otwarcia połączenia z bazą danych Postgres
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Kontrola, czy połączenie jest aktywne
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// W przypadku braku błędów zwróć wskaźnik połączenia db
	return db, nil
}

// Funkcja służy do nawiązywania połączenia z bazą danych
func connectToDB() *sql.DB {

	// Pobranie parametrów ze zmiennej środowiskowej DSN
	dsn := os.Getenv("DSN")

	// Kontynuuj dopóki Postgres nie będzie gotowy (limit prób 10)
	for {
		// Używamy funkcji openDB (zdefiniowanej powyżej)
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres nie jest jeszcze gotowy.")
			// Jeżeli się nie udało zwiększamy licznik prób połączeń
			counts++
		} else {
			log.Println("Połączony z Postgres!")
			// Jeżeli połączono, to funkcja zwraca "connection"
			return connection
		}

		if counts > 10 {
			// Po więcej niż 10 próbach, zwracamy nil
			log.Println(err)
			return nil
		}

		log.Println("Czekam dwie sekundy i sprawdzę ponownie.")
		time.Sleep(2 * time.Second)
		// Po dwóch sekundach ponawiamy próbę
		continue
	}

}

/* Rozbicie funkcji na dwie oddzielne funkcje pozwala na łatwiejsze
ponowne użycie kodu. Funkcja connectToDB() zawiera dodatkowe kroki
takie jak oczekiwanie na gotowość bazy danych. */
