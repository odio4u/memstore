# ---------- Build stage ----------
FROM golang:1.25 AS builder
WORKDIR /app

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build certificates
ARG ip
ARG dns
ARG name
RUN go run ./certmaker --ip "$ip" --dns "$dns" --name "$name" --force

# Prepare certs folder
RUN mkdir -p /app/bin/certs && \
    cp certmaker/server.pem /app/bin/certs/ && \
    cp certmaker/server-key.pem /app/bin/certs/

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/app cmd/main.go

# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /app/bin/app . 
COPY --from=builder /app/bin/certs ./certs

EXPOSE 50051
CMD ["/app/app"]
