# ---------- Build stage ----------
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source
COPY . .

# Build the application (same as: make build)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/app cmd/main.go

# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/app .

# Expose any port your app uses
EXPOSE 8080

# Run the binary
CMD ["/app/app"]
