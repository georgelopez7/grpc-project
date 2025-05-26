# 🧪 TESTS
test:
	go test ./...

# ---------------------------------------------------------------------------------------------------------------------------------------------------

# 🚀 START COMMANDS (FOR DEVELOPMENT)
start-gateway:
	go run cmd/gateway/main.go

start-fraud:
	go run cmd/fraud/main.go

# ---------------------------------------------------------------------------------------------------------------------------------------------------

# 📡 PROTOCOL GRPC GENERATION
genproto-payment:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/paymentpb/payment.proto

# ---------------------------------------------------------------------------------------------------------------------------------------------------

# 🐋 DOCKER COMPOSE
docker-up: # Builds and starts the docker-compose stack
	docker-compose up --build

docker-volume-clear: # Clears the docker volumes
	docker-compose down -v

AWS_REGION = eu-west-2
ACCOUNT_ID := $(shell aws sts get-caller-identity --query Account --output text)
REPO_NAME := gateway
IMAGE_TAG := latest
ECR_URI := $(ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(REPO_NAME)
DOCKER_PATH := docker/gateway.Dockerfile

push-docker-image-gateway:
	docker build -t $(REPO_NAME):$(IMAGE_TAG) -f $(DOCKER_PATH) .
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
	docker tag $(REPO_NAME):$(IMAGE_TAG) $(ECR_URI):$(IMAGE_TAG)
	docker push $(ECR_URI):$(IMAGE_TAG)
	
# ---------------------------------------------------------------------------------------------------------------------------------------------------

# 🧵 KUBERNETES & 🐢 HELM
ENV = dev
PATH_TO_HELM_K8s = infra/k8s/
PROJECT_NAME = grpc-project

# OS = windows| linux | mac
kube-deploy-local: # make kube-deploy-local OS=windows
	@echo "-> Starting Minikube..."
	minikube start

	@echo "-> Building Docker images inside Minikube (OS=$(OS))..."
ifeq ($(OS),windows)
	@powershell -Command "minikube docker-env | Invoke-Expression; \
		docker build -t gateway:latest -f docker/gateway.Dockerfile .; \
		docker build -t fraud:latest -f docker/fraud.Dockerfile .; \
		docker build -t validation:latest -f docker/validation.Dockerfile ."
else ifeq ($(OS),mac)
	@eval "$$(minikube docker-env)"; \
	docker build -t gateway:latest -f docker/gateway.Dockerfile .; \
	docker build -t fraud:latest -f docker/fraud.Dockerfile .; \
	docker build -t validation:latest -f docker/validation.Dockerfile .
else ifeq ($(OS),linux)
	@eval "$$(minikube docker-env)"; \
	docker build -t gateway:latest -f docker/gateway.Dockerfile .; \
	docker build -t fraud:latest -f docker/fraud.Dockerfile .; \
	docker build -t validation:latest -f docker/validation.Dockerfile .
else
	@echo "❌ Unsupported OS: $(OS). Please pass OS=windows, linux, or mac"
	@exit 1
endif

	@echo "-> Creating namespace '$(PROJECT_NAME)' if it doesn't exist..."
	kubectl create namespace "$(PROJECT_NAME)" --dry-run=client -o yaml | kubectl apply -f -

	@echo "-> Installing umbrella Helm chart and all dependencies ($(ENV) environment)..."
	helm upgrade --install $(PROJECT_NAME) $(PATH_TO_HELM_K8s) \
		-n $(PROJECT_NAME) \
		-f $(PATH_TO_HELM_K8s)/values.yaml

	@echo "-> Waiting for gateway deployment to be ready..."
	kubectl wait --for=condition=available --timeout=120s deployment/gateway -n $(PROJECT_NAME)

	@echo "-> Forwarding port to gateway service..."
	kubectl port-forward svc/gateway 8080:8080 -n $(PROJECT_NAME)
	@echo "-> Gateway available at http://localhost:8080/api/v1/payments"


kube-pods: # List all running pods in the cluster
	kubectl get pods -n $(PROJECT_NAME)

kube-services: # List all services in the cluster
	kubectl get services -n $(PROJECT_NAME)

kube-destroy: # Uninstalls Helm release and deletes namespace
	@echo "-> Uninstalling Helm release..."
	helm uninstall $(PROJECT_NAME) -n "$(PROJECT_NAME)" || true
	@echo "-> Deleting namespace..."
	kubectl delete namespace "$(PROJECT_NAME)" || true
	@echo "-> Stopping Minikube..."
	minikube stop
	@echo "-> Deleting Minikube cluster..."
	minikube delete

helm-lint: # Lints Helm charts
	helm lint $(PATH_TO_HELM_K8s)

helm-dependency-update: # Updates Helm dependencies
	helm dependency update $(PATH_TO_HELM_K8s)

# ---------------------------------------------------------------------------------------------------------------------------------------------------

# 👷🏻‍♂️ TERRAFORM (ECR)
PATH_TO_TERRAFORM_ECR = "infra/terraform/ecr"
terraform-init-ecr: # Initializes Terraform
	terraform -chdir=$(PATH_TO_TERRAFORM) init

terraform-lint-ecr: # Lints Terraform
	terraform -chdir=$(PATH_TO_TERRAFORM) validate

terraform-plan-ecr: # Dry run of Terraform
	terraform -chdir=$(PATH_TO_TERRAFORM) plan

terraform-apply-ecr: # Applies Terraform to AWS
	terraform -chdir=$(PATH_TO_TERRAFORM) apply

terraform-destroy-ecr:
	terraform -chdir=$(PATH_TO_TERRAFORM) destroy

# ---------------------------------------------------------------------------------------------------------------------------------------------------