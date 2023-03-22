package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go" //RMQ
)

// Stała określająca numer portu
const webPort = "80"

// RMQ: Licznik dla odliczania ilości prób połączeń z RabbitMQ
var counts int64

// Definicja struktury Config
type Config struct {
	Rabbit *amqp.Connection //RMQ
}

func main() {
	// Logowanie informacji o starcie usługi Broker'a na porcie 80
	log.Printf("Uruchamianie usługi Broker'a na porcie %s\n", webPort)

	// RMQ: Tworzenie połączenia do RabbitMQ
	rabbitConn, err := connectToRabbit()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// Inicjalizacja instancji aplikacji z konfiguracją Config
	app := Config{
		Rabbit: rabbitConn, //RMQ
	}

	// Definicja serwera HTTP
	srv := &http.Server{
		// Adres, na którym serwer będzie nasłuchiwać żądania HTTP
		Addr: fmt.Sprintf(":%s", webPort),
		// Funkcja obsługi żądań przychodzących na serwer
		Handler: app.routes(),
	}

	// Nasłuchiwanie na adresie i obsługa żądań HTTP
	err = srv.ListenAndServe()
	if err != nil {
		// W przypadku błędu zakończ program i zwróć błąd
		log.Panic(err)
	}
}

// RMQ:
// Funkcja służy do nawiązywania połączenia z RabbitMQ
// z wykorzystaniem biblioteki amqp, produkcji Streadway Labs
func connectToRabbit() (*amqp.Connection, error) {

	// Definicja wskaźnika do połączenia RabbitMQ
	var connection *amqp.Connection

	// Kontynuuj dopóki RabbitMQ nie będzie gotowy (limit prób 10)
	for {
		// Jeżeli listener-service mamy już w Dokerze, to łączymy się do "rabbitmq"
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		// Jeżeli sprawdzamy lokalnie "go run ." to odwołujemy się do "localhost"
		//c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			fmt.Println("RabbitMQ nie jest jeszcze gotowy.")
			// Jeżeli się nie udało zwiększamy licznik prób połączeń
			counts++
		} else {
			log.Println("Połączony z RabbitMQ!")
			// Jeżeli połączono, to funkcja zwraca "connection"
			connection = c
			return connection, nil
		}

		if counts > 10 {
			// Po więcej niż 10 próbach, zwracamy nil i błąd
			fmt.Println(err)
			return nil, err
		}

		log.Println("Czekam dwie sekundy i sprawdzę ponownie.")
		time.Sleep(2 * time.Second)
		// Po dwóch sekundach ponawiamy próbę
		continue
	}

}
