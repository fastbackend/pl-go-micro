# Definicja zmiennych
FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerApp
MAILER_BINARY=mailerApp
LISTENER_BINARY=listenerApp

#
# Sekcja BACKEND-u
#

# Uruchom wszystkie kontenery (bez wymuszania kompilacji)
up:
	@echo "Uruchom kontenery Dockera:"
	docker-compose up -d
	@echo "Kontenery Dockera uruchomione!"

# Zatrzymaj kontenery Dockera
down:
	@echo "Zatrzymaj kontenery Dockera:"
	docker-compose down
	@echo "Kontenery Dockera zatrzymane!"

#
# Sekcja przebudowy (kompilowania) BACKEND-u
#

# Kompiluj wszystkie projekty, zatrzymaj kontenery 
# i ponownie uruchom kontenery, uwzględniając aktualne parametry
up_build: build_broker build_auth build_logger build_mailer build_listener
	@echo "Zatrzymaj kontenery Dockera (jeżeli są uruchomione):"
	docker-compose down
	@echo "Kompiluj (jeżeli jest wymagane) i uruchom kontenery Dockera:"
	docker-compose up --build -d
	@echo "Kontenery Dockera uruchomione!"

# Kompiluj broker-service dla Linuxa
build_broker:
	@echo "Kompiluj broker-service:"
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Skompilowane!"

# Kompiluj auth-service dla Linuxa
build_auth:
	@echo "Kompiluj auth-service:"
	cd ../auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Skompilowane!"

# Kompiluj logger-service dla Linuxa
build_logger:
	@echo "Kompiluj logger-service:"
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Skompilowane!"

# Kompiluj mailer-service dla Linuxa
build_mailer:
	@echo "Kompiluj mailer-service:"
	cd ../mailer-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAILER_BINARY} ./cmd/api
	@echo "Skompilowane!"

# Kompiluj listener-service dla Linuxa
build_listener:
	@echo "Kompiluj listener-service:"
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Skompilowane!"

#
# Sekcja FRONTEND-u
#

# Uruchom frontend
start: build_front
	@echo "Uruchom Web Server (port 80):"
	cd ../frontend && go build -o ${FRONT_END_BINARY} ./cmd/web
	cd ../frontend && ./${FRONT_END_BINARY} &

# Zatrzymaj frontend
stop:
	@echo "Zatrzymaj Web Server (port 80):"
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Web Server zatrzymany!"

# Kompiluj frontend dla Linuxa
build_front:
	@echo "Kompiluj front-end:"
	cd ../frontend && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Skompilowane!"
