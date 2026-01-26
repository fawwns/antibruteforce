.PHONY: build run test clean down

APP_NAME=antibruteforce
BIN_DIR=bin
MAIN=./cmd/antibruteforce/main.go

build:
	@echo "🔨 Building $(APP_NAME)..."
	go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN)

run:
	@echo "🐳 Starting service with docker-compose..."
	docker compose up --build

down:
	@echo "🛑 Stopping docker-compose..."
	docker compose down

test:
	@echo "🧪 Running tests..."
	go test ./...

clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BIN_DIR)