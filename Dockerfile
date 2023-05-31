FROM golang:alpine3.18

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o app cmd/main.go

CMD ["./app"]
