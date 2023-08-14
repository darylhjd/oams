# syntax=docker/dockerfile:1
# Context: Project root.

# Stage 1: Build flutter web.
FROM ghcr.io/cirruslabs/flutter:3.10.6 as flutter-builder

WORKDIR /app
COPY ./frontend ./

RUN --mount=type=secret,id=env_json flutter pub get && flutter build web --dart-define-from-file=/run/secrets/env_json

# Stage 2: Build bundledserver.
FROM golang:1.20-alpine as builder

RUN apk add --no-cache git \
    openssh-client \
    ca-certificates

WORKDIR /app

ENV GOPRIVATE="github.com/darylhjd/oams/*"
ENV GIT_SSH_COMMAND="ssh -i /run/secrets/ssh_key"
RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"
RUN mkdir -p -m 0700 /root/.ssh && ssh-keyscan github.com >> /root/.ssh/known_hosts
COPY ./backend/go.mod ./backend/go.sum ./
RUN --mount=type=secret,id=ssh_key go mod download

COPY ./backend ./
# Copy flutter build.
COPY --from=flutter-builder /app/build/web/ ./internal/servers/webserver/build/web/

WORKDIR /app/cmd/bundledserver
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out

# Stage 3: Move bundledserver binary to bare container.
FROM scratch
COPY --from=builder /out /out

# Required for making HTTP requests in the container.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["/out"]