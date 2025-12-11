.PHONY: all clean build-windows build-darwin build-linux build-all test

# Default target
all: build-all

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o pos-printer-windows.exe

# Build for macOS (both architectures)
build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o pos-printer-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o pos-printer-darwin-arm64

# Build for Linux (both architectures)
build-linux:
	GOOS=linux GOARCH=amd64 go build -o pos-printer-linux-amd64
	GOOS=linux GOARCH=arm64 go build -o pos-printer-linux-arm64

# Build for all platforms
build-all: build-windows build-darwin build-linux
	@echo "All builds completed successfully!"

# Clean build artifacts
clean:
	rm -f pos-printer-windows.exe
	rm -f pos-printer-darwin-amd64
	rm -f pos-printer-darwin-arm64
	rm -f pos-printer-linux-amd64
	rm -f pos-printer-linux-arm64
	rm -f pos-printer.exe

# Run tests
test:
	go test ./...

# Run the application (native build)
run:
	go run main.go
