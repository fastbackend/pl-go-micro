# Tworzymy nowy obraz na bazie alpine:latest
FROM alpine:latest

# Tworzymy katalog /app
RUN mkdir /app

# Kopiujemy plik binarny authApp do /app w tworzonym obrazie
# oraz szablony dla wiadomości HTML do templates
COPY mailerApp /app
COPY templates /templates

# Uruchamiamy aplikację, kiedy kontener jest uruchamiany
CMD [ "/app/mailerApp" ]