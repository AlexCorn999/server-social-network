FROM golang:latest
WORKDIR /app
COPY . /app

# Build the go app
RUN go build -o ./cmd/main ./cmd


# Expose port
EXPOSE 8080

# Define the command to run the app
ENTRYPOINT ["./cmd/main"]
