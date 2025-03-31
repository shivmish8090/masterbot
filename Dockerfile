FROM golang:1.23-alpine

ENV GOPATH=/app/go
ENV PATH=$GOPATH/bin:$PATH

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

# Copy module sources from `pkg/mod` to `src/` (needed for Yaegi)
RUN mkdir -p $GOPATH/src && cp -r $GOPATH/pkg/mod/* $GOPATH/src/

COPY . .

RUN go build -o app .

CMD ["./app"]