FROM golang:1.23-alpine

ENV GOPATH=/app/go
ENV PATH=$GOPATH/bin:$PATH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN go build -o app .

CMD ["./app"]