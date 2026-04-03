# ---- Build Stage ----
FROM golang:1.26-alpine AS builder
ARG PROJECT_DIR

# Set working directory
WORKDIR /build

COPY lib ./lib
COPY gen ./gen
COPY services ./services
COPY go.work go.work.sum ./

# Build the executable for Linux, statically linked
# The main package is in ./server
RUN go build -ldflags="-w -s" \
    -o server ./${PROJECT_DIR}/server

# ---- Final Stage ----
FROM alpine:latest
ARG PROJECT_DIR

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Create application directory and set ownership
WORKDIR /app
RUN chown appuser:appgroup /app

# Copy the binary from builder with correct ownership
COPY --from=builder --chown=appuser:appgroup /build/server .

# Switch to non‑root user
USER appuser

# Run the server
ENTRYPOINT ["./server"]