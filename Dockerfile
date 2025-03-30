#FROM golang:1.23-bookworm
FROM golang:1.23-alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN go build -o app .

CMD ["./app"]