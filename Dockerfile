FROM golang:1.19.0

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o cmd/main cmd/main.go

EXPOSE 9000

CMD ["./cmd/main"]