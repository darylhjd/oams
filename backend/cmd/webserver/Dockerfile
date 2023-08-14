# syntax=docker/dockerfile:1
FROM ghcr.io/cirruslabs/flutter:3.10.6 as flutter-builder

WORKDIR /app
COPY ./frontend ./

RUN flutter pub get && flutter build web

FROM golang:1.20-alpine as builder

WORKDIR /app
COPY ./backend ./
RUN go mod download

# Copy the project local code.
COPY --from=flutter-builder /app/build/web/ ./internal/servers/webserver/build/web/

WORKDIR /app/cmd/webserver
RUN CGO_ENABLED=0 GOOS=linux go build -o /out

FROM scratch
COPY --from=builder /out /out

# Required for making HTTP requests in the container.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["/out"]