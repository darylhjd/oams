# syntax=docker/dockerfile:1

FROM golang:1.20-alpine as builder

ARG pat
ENV github_pat $pat

RUN apk update && apk add git

# Set working directory.
WORKDIR /app

# Copy go module files and get dependencies. Note that we are running this from project root!
COPY go.* ./
RUN echo "$github_pat"
RUN git config --global --add url."https://darylhjd:$github_pat@github.com".insteadOf "https://github.com"  \
    && GOPRIVATE="github.com/darylhjd/*"
RUN go mod download

# Copy the project local code.
COPY . ./

WORKDIR /app/cmd/apiserver
RUN CGO_ENABLED=0 GOOS=linux go build -o /out

# Add a new user "nonroot".
RUN addgroup -S nonrootgroup && adduser -S nonroot -G nonrootgroup

FROM scratch
COPY --from=builder /out /out

# Change to non-root privilege
COPY --from=builder /etc/passwd /etc/passwd
USER nonroot

EXPOSE 8080

ENTRYPOINT ["/out"]