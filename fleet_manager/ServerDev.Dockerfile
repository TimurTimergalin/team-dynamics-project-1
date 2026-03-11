# ---- Build Stage ----
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /build

# Copy go module files first (for better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source
COPY . .

# Build the executable for Linux, statically linked
# The main package is in ./server
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-w -s" \
    -o fleet-manager-server ./server

# ---- Final Stage ----
FROM alpine:latest

# Create application directory and set ownership
WORKDIR /app
RUN chown appuser:appgroup /app

# Copy the binary from builder with correct ownership
COPY --from=builder --chown=appuser:appgroup /build/fleet-manager-server .

# Switch to non‑root user
USER appuser

# Run the server
ENTRYPOINT ["./fleet-manager-server"]