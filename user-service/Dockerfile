# Build the app
FROM golang:1.19.0-alpine3.16 AS builder

WORKDIR /backend-user-service

COPY user-service/go.mod ./
COPY user-service/go.sum ./

RUN mkdir ../backend-common
COPY backend-common ../backend-common

RUN mkdir src
RUN mkdir bin

RUN go mod download

COPY user-service/src/*.go src/

RUN go build -o ./bin/user-service ./src

# Copy it the app to an alpine
FROM alpine:3.11.3

WORKDIR /app

COPY --from=builder /backend-user-service/bin/user-service ./

EXPOSE 8000

HEALTHCHECK  --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/healthcheck || exit 1

CMD ["./user-service"]
