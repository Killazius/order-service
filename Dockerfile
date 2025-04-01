
FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o order-server ./cmd/order-server/
CMD ["./order-server"]

EXPOSE 8080
EXPOSE 50051