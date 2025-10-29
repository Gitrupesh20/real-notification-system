# ---- Build Stage ----
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first (to leverage caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary (optimized for production)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app -trimpath -ldflags="-s -w" .

# ---- Runtime Stage ----
FROM alpine:latest AS runner

# For minimal image size & security
RUN adduser -D appuser

WORKDIR /myBin

# Copy only the binary from builder
COPY --from=builder /app/bin/app .

# Set permissions
RUN chown appuser:appuser /myBin/app
USER appuser

EXPOSE 5050

# Run the binary
CMD ["./app"]
