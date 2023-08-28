test:
	go test -v -cover ./...

server:
	go run main.go

proto:
	rm internal/pb/*
	protoc --go_out=internal/pb --go_opt=paths=source_relative --go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=internal/pb --grpc-gateway_opt=paths=source_relative --proto_path=internal/proto internal/proto/*.proto

.PHONY: test server proto