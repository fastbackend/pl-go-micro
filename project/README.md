To jest przykład materiału do szkolenia z mikroserwisów w Golangu, z wykorzystaniem RabbitMQ
(komentarze są po polsku, ponieważ materiał jest przeznaczony do szkoleń polskojęzycznych)

Notatki pomocne do uruchamiania projektu pl-go-micro:

Wymagania wstępne:
1. Golang 1.20 (dla kompilowania kodu)
2. Docker Desktop (na potrzeby uruchamiania mikroserwisów i baz danych)
3. Visual Studio Code (jeżeli chcesz modyfikować kod)
(projekt realizowano w środowisku MacBook/Linux z procesorem arm)

W projekcie używam Workspace w ramach Visual Studio Code, który nie pozwala dołączyć pliku jako Workspace, ale zawsze to musi być folder, co wymusza docker-compose.yml z podaniem contextu i ścieżki do dockerfile'a.

Zależności w obszarze broker-service:
- go get github.com/go-chi/chi/v5
- go get github.com/go-chi/chi/v5/middleware
- go get github.com/go-chi/cors

i dodatkowo po włączeniu korzystania z RabbitMQ
- go get github.com/rabbitmq/amqp091-go

Zależności w obszarze auth-service:
- te co w broker-service (*go-chi*) plus:
- go get github.com/jackc/pgconn
- go get github.com/jackc/pgx/v4
- go get github.com/jackc/pgx/v4/stdlib

Hasło testowe dla hubert@example.com: verysecret
(przy pierwszym utworzeniu postgres wykonać skrypt users.sql z katalogu project/support)

Zależności w obszarze logger-service:
- te co w broker-service (*go-chi*) plus:
- go get go.mongodb.org/mongo-driver/mongo
- go get go.mongodb.org/mongo-driver/mongo/options
- go get go.mongodb.org/mongo-driver/bson

Parametry połączenia do bazy mongoDB (np. poprzez MongoDB Compass): 
mongodb://admin:password@localhost:27017/logs?authSource=admin&readPreference=primary&appname=MongDB%20Compass&directConnection=true&ssl=false

Zależności w obszarze mailer-service:
- te co w broker-service (*go-chi*) plus:
- go get github.com/vanng822/go-premailer/premailer
- go get github.com/xhit/go-simple-mail/v2

Zależności w obszarze listener-service (RabbitMQ), amqp:
- go get github.com/rabbitmq/amqp091-go
