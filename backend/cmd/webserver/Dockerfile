# syntax=docker/dockerfile:1
# Context: Project root.

# Stage 1: Build web.
FROM node:lts-alpine3.18 as builder

WORKDIR /app
COPY ./frontend ./

RUN --mount=type=secret,id=env cat /run/secrets/env > /app/.env.production  \
    && npx next telemetry disable  \
    && npm install  \
    && npm run build

# Stage 2: New container to run.
FROM node:lts-alpine3.18 as runner

WORKDIR ./app

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Required for making HTTP requests in the container.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD ["node", "server.js"]