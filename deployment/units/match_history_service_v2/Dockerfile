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
ARG PROJECT_DIR

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

#RUN apk add --no-cache \
#    net-tools \
#    bind-tools \
#    busybox-extras \
#    curl \
#    wget \
#    procps \
#    htop \
#    tree \
#    vim \
#    less \
#    postgresql-client \
#    bash \
#    jq \
#    mtr \
#    traceroute

WORKDIR /app
RUN chown appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /build/server .

USER appuser

ENTRYPOINT ["./server"]
