BROKER_BINARY=brokerApp
SEARCH_BINARY=searchServiceApp
MED_SCRAPER_BINARY=medScraperServiceApp
MED_API_BINARY=medApiServiceApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_search build_pubmed build_wiki
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_search: builds the search binary as a linux executable
build_search:
	@echo "Building search binary..."
	cd ../search-service && env GOOS=linux CGO_ENABLED=0 go build -o ${SEARCH_BINARY} ./cmd/api
	@echo "Done!"

build_pubmed:
	@echo "Building med scraper binary..."
	cd ../med-scraper-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MED_SCRAPER_BINARY} ./cmd/api
	@echo "Done!"

build_wiki:
	@echo "Building med api service binary..."
	cd ../med-api-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MED_API_BINARY} ./cmd/api
	@echo "Done!"
