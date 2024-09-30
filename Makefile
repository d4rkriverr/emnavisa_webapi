# Variables
APP_NAME := ./bin/myapp
SRC := ./cmd/main.go

# Build the application
build:
	@echo "Building the Go application..."
	@go build -o $(APP_NAME) $(SRC)

# Run the application
run: build
	@echo "Running the Go application..."
	@./$(APP_NAME)
