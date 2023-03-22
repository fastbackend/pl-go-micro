# v1: Wersja 1 była wykorzystywana, gdy nie używałem Makefile
# zostawiłem to jako przykład budowania aplikacji z kodu w kontenerze Dockera

# Ustawiamy obraz bazowy na golang:1.20-alpine
FROM golang:1.20-alpine as builder

# Tworzymy katalog /app
RUN mkdir /app

# Kopiujemy zawartość bieżącego katalogu do /app
# a bieżący katalog to broker-service
COPY . /app

# Ustawiamy katalog roboczy na /app
WORKDIR /app

# Budujemy aplikację Go, wyłączamy CGO i tworzymy wykonywalny plik binarny o nazwie brokerApp w katalogu /app
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# Nadajemy uprawnienia do wykonania pliku binarnego
RUN chmod +x /app/brokerApp

# Tworzymy nowy obraz na bazie alpine:latest
FROM alpine:latest

# Tworzymy katalog /app
RUN mkdir /app

# Kopiujemy plik binarny brokerApp z obrazu budującego do /app w nowym obrazie
COPY --from=builder /app/brokerApp /app

# Uruchamiamy aplikację, kiedy kontener jest uruchamiany
CMD [ "/app/brokerApp" ]