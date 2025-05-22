# 🧪 TESTS
test:
	go test ./...



# 🚀 START COMMANDS (FOR DEVELOPMENT)
start-gateway:
	go run cmd/gateway/main.go

start-fraud:
	go run cmd/fraud/main.go



# 📡 PROTOCOL GRPC GENERATION
genproto-payment:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/paymentpb/payment.proto



# 🐋 DOCKER COMPOSE
docker-up: # Builds and starts the docker-compose stack
	docker-compose up --build

docker-volume-clear: # Clears the docker volumes
	docker-compose down -v