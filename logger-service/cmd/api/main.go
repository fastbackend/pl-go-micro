package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fastbackend/pl-go-micro/logger-service/data"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Stałe określające parametry łączenia z mongoDB
const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

// Wskaźnik obiektu klienta mongoDB
var client *mongo.Client

// Definicja struktury Config
type Config struct {
	Models data.Models
}

func main() {

	// Logowanie informacji o starcie usługi autentykacji na porcie 80
	log.Printf("Uruchamianie usługi dziennika zdarzeń (logger'a) na porcie %s\n", webPort)

	// Tworzenie połączenia do bazy danych
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic("Nie można połączyć się z mongoDB!")
	}
	client = mongoClient

	// Utwórz kontekst reguły rozłączenia
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Zamknij połączenie
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Inicjalizacja instancji aplikacji z konfiguracją Config
	app := Config{
		Models: data.New(client),
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

// Funkcja służy do nawiązywania połączenia z bazą danych
func connectToMongo() (*mongo.Client, error) {

	// Tworzenie parametrów dla połączenia
	// (wyrzucić te parametry do zmiennych środowiskowych)
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// Tworzenie połączenia z bazą danych
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Błąd połączenia z mongoDB:", err)
		return nil, err
	}

	log.Println("Połączony z mongoDB!")

	return c, nil
}
