# ---------- Build stage ----------
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source
COPY . .

ARG ip
ARG dns
ARG name

COPY certmaker ./certmaker

RUN go run ./certmaker --ip "$ip" --dns "$dns" --name "$name" --force

RUN mkdir -p /app/bin/certs && \
    cp certmaker/server.pem /app/bin/certs/ && \
    cp certmaker/server-key.pem /app/bin/certs/


# Build the application (same as: make build)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/app cmd/main.go

# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/app .
COPY --from=builder /app/bin/certs ./certs


# Expose any port your app uses
EXPOSE 50051

# Run the binary
CMD ["/app/app"]
