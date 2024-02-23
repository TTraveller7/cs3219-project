FROM golang:1.19.0-alpine3.16 AS builder

WORKDIR /app/matching-service

COPY backend-matching-service/matching-service/go.mod ./
COPY backend-matching-service/matching-service/go.sum ./

RUN mkdir /backend-common
COPY backend-common /backend-common

RUN mkdir {bin,cmd,match}

RUN go mod download

COPY backend-matching-service/matching-service/cmd ./cmd
COPY backend-matching-service/matching-service/match ./match

RUN go build -o ./bin/matching-service ./cmd

# Copy it the app to an alpine
FROM alpine:3.11.3

WORKDIR /app

COPY --from=builder /app/matching-service/bin/matching-service ./

EXPOSE 7001

HEALTHCHECK  --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:7001/healthcheck || exit 1

CMD ["./matching-service"]
