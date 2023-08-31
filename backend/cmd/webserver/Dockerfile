# syntax=docker/dockerfile:1
# Context: Project root.

# Stage 1: Build web.
FROM node:lts-alpine3.18 as web-builder

WORKDIR /app
COPY ./frontend ./

RUN --mount=type=secret,id=env cat /run/secrets/env > /app/.env.production  \
    && npx next telemetry disable  \
    && npm install  \
    && npm run build

# Stage 2: Build webserver.
FROM golang:1.20-alpine as builder

WORKDIR /app
COPY ./backend ./
RUN go mod download

# Copy flutter build.
COPY --from=web-builder /app/out/ ./internal/servers/webserver/build/out/

WORKDIR /app/cmd/webserver
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out

# Stage 3: Move webserver binary to bare container.
FROM scratch
COPY --from=builder /out /out

# Required for making HTTP requests in the container.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["/out"]