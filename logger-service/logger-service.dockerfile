# Tworzymy nowy obraz na bazie alpine:latest
FROM alpine:latest

# Tworzymy katalog /app
RUN mkdir /app

# Kopiujemy plik binarny authApp do /app w tworzonym obrazie
COPY loggerApp /app

# Uruchamiamy aplikację, kiedy kontener jest uruchamiany
CMD [ "/app/loggerApp" ]