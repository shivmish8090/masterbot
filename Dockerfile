FROM golang:1.22-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o app .

FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

COPY --from=builder /app/app .

RUN chmod +x ./app

CMD ["./app"]