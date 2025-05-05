test:
	go test ./...

start-gateway:
	go run cmd/gateway/main.go

start-fraud:
	go run cmd/fraud/main.go

genproto-payment:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/paymentpb/payment.proto