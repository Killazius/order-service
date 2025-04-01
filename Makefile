BINARY_NAME=order-server

build:
	go build -o ${BINARY_NAME} cmd/order-server/order-server.go

exec:
	./${BINARY_NAME}

clear:
	del ${BINARY_NAME}
run:
	make build && make exec && make clear
proto:
	protoc --go_out=./pkg/api/test --go_opt=paths=source_relative \
        --go-grpc_out=./pkg/api/test --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=./pkg/api/test --grpc-gateway_opt=paths=source_relative \
        order.proto
