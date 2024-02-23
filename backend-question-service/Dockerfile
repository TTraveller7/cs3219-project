FROM golang:1.19.0-alpine3.16 AS builder

WORKDIR /backend-question-service

COPY backend-question-service/go.mod ./
COPY backend-question-service/go.sum ./

RUN mkdir /backend-common
COPY backend-common /backend-common

RUN mkdir bin

RUN go mod download

COPY backend-question-service/*.go ./

RUN go build -o ./bin/question-service ./

# Copy it the app to an alpine
FROM alpine:3.11.3

WORKDIR /app/bin

COPY --from=builder /backend-question-service/bin/question-service ./

EXPOSE 17001

HEALTHCHECK  --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:17001/healthcheck || exit 1

CMD ["./question-service"]
