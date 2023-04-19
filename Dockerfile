FROM golang:1.20-alpine

WORKDIR /app

RUN export GO111MODULE=on
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/ethereum_parser

RUN go build -o /

EXPOSE 8080

CMD ["/ethereum_parser"]
