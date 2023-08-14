# syntax=docker/dockerfile:1
FROM golang:1.20-alpine as builder

RUN apk add --no-cache git \
    openssh-client \
    ca-certificates

# Set working directory.
WORKDIR /app

# Copy go module files and get dependencies. Note that we are running this from project root!
ENV GOPRIVATE="github.com/darylhjd/oams/*"
ENV GIT_SSH_COMMAND="ssh -i /run/secrets/ssh_key"
RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"
RUN mkdir -p -m 0700 /root/.ssh && ssh-keyscan github.com >> /root/.ssh/known_hosts
COPY go.mod go.sum ./
RUN --mount=type=secret,id=ssh_key go mod download

# Copy the project local code.
COPY . ./

WORKDIR /app/cmd/apiserver
RUN CGO_ENABLED=0 GOOS=linux go build -o /out

FROM scratch
COPY --from=builder /out /out

# Required for making HTTP requests in the container.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["/out"]