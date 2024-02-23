FROM golang:1.19.0-alpine3.16 AS builder

WORKDIR /backend-auth-service

COPY backend-auth-service/go.mod ./
COPY backend-auth-service/go.sum ./

RUN mkdir /backend-common
COPY backend-common /backend-common

RUN mkdir bin

RUN go mod download

COPY backend-auth-service/*.go ./

RUN go build -o ./bin/auth-service ./


# Copy it the app to an alpine
FROM alpine:3.11.3

WORKDIR /app

COPY --from=builder /backend-auth-service/bin/auth-service ./

EXPOSE 13704

HEALTHCHECK  --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:13704/healthcheck || exit 1

CMD ["./auth-service"]
