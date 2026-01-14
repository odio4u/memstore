### Application commands

help:
	@echo "Available commands:"
	@echo "  make install-deps   - Install project dependencies"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make proto-gen      - Generate protobuf code"
	@echo "  make wal-proto-gen  - Generate WAL protobuf code"
	@echo "  make gen-cert       - Generate TLS certificates"

install-deps: .env
	go mod tidy

build: .env
	go build -o bin/app cmd/main.go

run: .env
	go run cmd/main.go

proto-gen:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/maps.proto

wal-proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative wal/proto/wal.proto

gen-cert:
	@echo "Generating certificates..."
	cd certmaker && \
	go run . \
	  --ip "$(ip)" \
	  --dns "$(dns)" \
	  --name "$(name)" \
	  $(if $(filter 1,$(client)),--client,) \
	  $(if $(filter 1,$(force)),--force,) \
	;


.PHONY: install-deps build run proto-gen wal-proto-gen gen-cert help