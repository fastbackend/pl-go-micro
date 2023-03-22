# v2: Wersja 2 bazuje na już skompilowanej aplikacji Brokera
# warunkiem jest, że mamy taką skomilowaną aplikację (z użyciem Makefile)

# Tworzymy nowy obraz na bazie alpine:latest
FROM alpine:latest

# Tworzymy katalog /app
RUN mkdir /app

# Kopiujemy plik binarny brokerApp do /app w tworzonym obrazie
COPY brokerApp /app

# Uruchamiamy aplikację, kiedy kontener jest uruchamiany
CMD [ "/app/brokerApp" ]