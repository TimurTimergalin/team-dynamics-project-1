# ---- Build Stage ----
FROM golang:1.26-alpine AS builder

# Set working directory
WORKDIR /build

# Copy go module files first (for better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source
COPY . .

# Build the executable for Linux, statically linked
# The main package is in ./client
RUN go build -ldflags="-w -s" \
    -o simple-test-client ./client

# ---- Final Stage ----
FROM alpine:latest

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Create application directory and set ownership
WORKDIR /app
RUN chown appuser:appgroup /app

# Copy the binary from builder with correct ownership
COPY --from=builder --chown=appuser:appgroup /build/simple-test-client .

# Switch to non‑root user
USER appuser

# Run the server
ENTRYPOINT ["./simple-test-client"]