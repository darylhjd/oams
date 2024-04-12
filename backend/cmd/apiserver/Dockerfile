# syntax=docker/dockerfile:1
# Context: Backend.

# Stage 1: Build apiserver.
FROM golang:1.22-alpine as builder

RUN apk add --no-cache git \
    openssh-client \
    ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

WORKDIR /app/cmd/apiserver
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out

# Stage 2: Move apiserver binary to bare container.
FROM scratch
COPY --from=builder /out /out

# Required for making HTTP requests in the container.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/out"]