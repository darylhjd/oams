# syntax=docker/dockerfile:1

FROM golang:1.20

ARG pat
ARG port
ENV github_pat $pat
ENV apiserver_port $port

WORKDIR /app

RUN git config --global --add url."https://darylhjd:$github_pat@github.com".insteadOf "https://github.com"  \
    && GOPRIVATE="github.com/darylhjd/*"

COPY main.go ./
RUN go mod init && go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build /apiserver

EXPOSE $apiserver_port

CMD ["/apiserver"]
