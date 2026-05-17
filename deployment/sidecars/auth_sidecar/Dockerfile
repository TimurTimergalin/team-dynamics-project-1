FROM golang:1.26-alpine AS builder
ARG PROJECT_DIR

WORKDIR /build

COPY lib ./lib
COPY gen ./gen
COPY services ./services
COPY go.work go.work.sum ./

RUN go build -ldflags="-w -s" \
    -o server ./${PROJECT_DIR}/server

FROM alpine:latest

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app
RUN chown appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /build/server .

USER appuser

ENTRYPOINT ["./server"]
