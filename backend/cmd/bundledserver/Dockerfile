# syntax=docker/dockerfile:1
# Context: Project root.

# Stage 1: Build web.
FROM node:lts-alpine3.18 as web-builder

WORKDIR /app
COPY ./frontend ./

RUN --mount=type=secret,id=env cat /run/secrets/env > /app/.env.production  \
    && npm next telemetry disable  \
    && npm install  \
    && npm run build

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
# Copy web build.
COPY --from=web-builder /app/out/ ./internal/servers/webserver/build/out/

WORKDIR /app/cmd/bundledserver
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out

# Stage 3: Move bundledserver binary to bare container.
FROM scratch
COPY --from=builder /out /out

# Required for making HTTP requests in the container.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["/out"]