# syntax=docker/dockerfile:1
# Context: Project root.

# Stage 1: Build flutter web.
FROM ghcr.io/cirruslabs/flutter:3.13.0 as flutter-builder

WORKDIR /app
COPY ./frontend ./

RUN --mount=type=secret,id=env_json flutter pub get && flutter build web --dart-define-from-file=/run/secrets/env_json

# Stage 2: Build webserver.
FROM golang:1.20-alpine as builder

WORKDIR /app
COPY ./backend ./
RUN go mod download

# Copy flutter build.
COPY --from=flutter-builder /app/build/web/ ./internal/servers/webserver/build/web/

WORKDIR /app/cmd/webserver
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out

# Stage 3: Move webserver binary to bare container.
FROM scratch
COPY --from=builder /out /out

# Required for making HTTP requests in the container.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["/out"]