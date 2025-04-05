FROM golang:1.24-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o app .

FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/app /

ENTRYPOINT ["/app"]