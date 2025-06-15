FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y libc6

WORKDIR /app

COPY --from=builder /app/main .

COPY .env /app/.env
COPY database/schema.sql /app/database/schema.sql
COPY mqtt/track.csv /app/mqtt/track.csv

EXPOSE 8080

CMD ["./main"]
