.PHONY: build clean run install test

clean :
	@echo "Cleaning builds..."
	@rm -rf build/server
	@echo "Clean complete."
	
build: clean
	@echo "Building the application..."
	@go build -o build/server .
	@echo "Build complete."


run: build
	@./build/server

install:
	@echo "Installing dependencies..."
	@go mod tidy
	@echo "Dependencies installed."

test:
	@go test ./...