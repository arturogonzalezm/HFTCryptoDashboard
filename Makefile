# Makefile for RealTimeBinanceMonitor

.PHONY: all build up down clean logs init-db

# Default target
all: build up

rebuild: clean build up

# Build the Docker images without using cache
build:
	@echo "Building Docker images..."
	docker compose build --no-cache

# Start the Docker containers in detached mode
up:
	@echo "Starting Docker containers..."
	docker compose up -d

# Stop the Docker containers
down:
	@echo "Stopping Docker containers..."
	docker compose down

# Remove Docker volumes
clean:
	@echo "Removing Docker volumes..."
	docker compose down -v
	docker volume rm timescaledb || true

# Show logs for all containers
logs:
	@echo "Showing logs for all containers..."
	docker compose logs -f

db-shell:
	docker exec -it timescaledb psql -U postgres -d hft


#run:
#	go run backend/main.go

#run:
#	@echo "Running the Go application..."
#	DB_HOST=localhost DB_USER=postgres DB_PASSWORD=postgres DB_NAME=hft go run backend/main.go

# Run the Go application (for local development)
run:
	@echo "Running the Go application..."
	env $$(cat .env | xargs) go run backend/main.go