### Application commands

help:
	@echo "Available commands:"
	@echo "  make install-deps      - Install project dependencies"
	@echo "  make build             - Build the application for current OS"
	@echo "  make build-all         - Build the application for all platforms"
	@echo "  make run               - Run the application"
	@echo "  make proto-gen         - Generate protobuf code"
	@echo "  make wal-proto-gen     - Generate WAL protobuf code"
	@echo "  make gen-cert          - Generate TLS certificates"

install-deps: .env
	go mod tidy

# Build for the current OS
build: .env
	go build -o bin/seeder cmd/main.go

# Cross-compile for all platforms
build-all: .env
	@echo "Building for Windows, Linux, and macOS..."
	@mkdir -p release
	@echo "Building Linux..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o release/seeder-linux-amd64 cmd/main.go
	@echo "Building macOS..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o release/seeder-macos-amd64 cmd/main.go
	@echo "Building Windows..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o release/seeder-windows-amd64.exe cmd/main.go
	@echo "All builds completed."

run: .env
	go run cmd/main.go

proto-gen:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/maps.proto

wal-proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative wal/proto/wal.proto

gen-cert:
	@echo "Generating certificates..."
	go run certs/seedcert.go

.PHONY: install-deps build build-all run proto-gen wal-proto-gen gen-cert help