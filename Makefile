AUTH_BINARY=authApp
BROKER_BINARY=brokerApp
FRONT_END_BINARY=frontendApp
LISTENER_BINARY=listenerApp
LOGGER_BINARY=loggerApp
MAIL_BINARY=mailApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_mail build_listener
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o bin/${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -o bin/${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd ./listener-service && env GOOS=linux CGO_ENABLED=0 go build -o bin/${LISTENER_BINARY} .
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_mail:
	@echo "Building mail binary..."
	cd ./mail-service && env GOOS=linux CGO_ENABLED=0 go build -o bin/${MAIL_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ./auth-service && env GOOS=linux CGO_ENABLED=0 go build -o bin/${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_front: builds the frone end binary
build_front_linux:
	@echo "Building front end linux binary..."
	cd ./front-end && env GOOS=linux CGO_ENABLED=0 go build -o bin/${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## build_front: builds the frone end binary
build_front_mac:
	@echo "Building front end mac binary..."
	cd ./front-end && env CGO_ENABLED=0 go build -o bin/${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front_mac
	@echo "Starting front end"
	cd ./front-end && ./bin/${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./bin/${FRONT_END_BINARY}"
	@echo "Stopped front end!"
