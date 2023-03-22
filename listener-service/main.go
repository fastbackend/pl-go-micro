package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fastbackend/pl-go-micro/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Licznik dla odliczania ilości prób połączeń z RabbitMQ
var counts int64

func main() {
	// Wyświetlanie informacji o rozpoczęciu usługi kolejkowania
	log.Println("Uruchamianie usługi kolejkowania")

	// Tworzenie połączenia do RabbitMQ
	rabbitConn, err := connectToRabbit()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// Wyświetlanie informacji o powodzeniu
	log.Println("Połączony z RabbitMQ!")

	// Informacja o uruchomieniu konsumera
	log.Println("Nasłuchiwanie i konsumowanie wiadomości przez RabbitMQ")

	// Tworzenie konsumera
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// Kolejkowanie i konsumowanie zdarzeń
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}

}

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
