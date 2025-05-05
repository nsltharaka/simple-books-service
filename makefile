build:
	@echo "Building the application..."
	@go build -o build/server .
	@echo "Build complete."

clean :
	@echo "Cleaning builds..."
	@rm -rf build/server
	@echo "Clean complete."

run: build
	@./build/server

install:
	@echo "Installing dependencies..."
	@go mod tidy
	@echo "Dependencies installed."

test:
	@go test ./...