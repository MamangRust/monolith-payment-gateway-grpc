set shell := ["bash", "-c"]

COMPOSE_FILE := "deployments/local/docker-compose.yml"
SERVICES := "apigateway migrate auth user role card merchant saldo topup transaction transfer withdraw email"
DOCKER_COMPOSE := "docker compose"
PROTO_DIR := "proto"
OUTDIR_PROTO := "pb"

# List all recipes
default:
    @just --list

# Run database migrations up
migrate:
    go run service/migrate/main.go up

# Run database migrations down
migrate-down:
    go run service/migrate/main.go down

# Generate protocol buffers
generate-proto:
    protoc \
        --proto_path={{PROTO_DIR}} \
        --go_out={{OUTDIR_PROTO}} --go_opt=paths=source_relative \
        --go-grpc_out={{OUTDIR_PROTO}} --go-grpc_opt=paths=source_relative \
        $(find {{PROTO_DIR}} -name "*.proto")

# Generate sqlc code
generate-sql:
    sqlc generate

# Generate swagger documentation
generate-swagger:
    swag init -g service/apigateway/cmd/main.go -o service/apigateway/docs

# Run seeder
seeder:
    go run service/seeder/main.go

# Build docker images for all services
build-image:
    @for service in {{SERVICES}}; do \
        echo "🔨 Building $service-service..."; \
        docker build -t $service-service:1.1 -f service/$service/Dockerfile service/$service || exit 1; \
    done
    @echo "✅ All services built successfully."

# Load images to minikube
image-load:
    @for service in {{SERVICES}}; do \
        echo "🚚 Loading $service-service..."; \
        minikube image load $service-service:1.1 || exit 1; \
    done
    @echo "✅ All services loaded successfully."

# Delete images from minikube
image-delete:
    @for service in {{SERVICES}}; do \
        echo "🗑️ Deleting $service-service image..."; \
        minikube image rm $service-service:1.1 || echo "⚠️ Failed to delete $service-service (maybe not found)"; \
    done
    @echo "✅ All requested images deleted (if they existed)."

# Show docker compose process status
ps:
    {{DOCKER_COMPOSE}} -f {{COMPOSE_FILE}} ps

# Start docker compose services
up:
    {{DOCKER_COMPOSE}} -f {{COMPOSE_FILE}} up -d

# Stop docker compose services
down:
    {{DOCKER_COMPOSE}} -f {{COMPOSE_FILE}} down

# Build images and start docker compose
build-up: build-image up

# Start minikube with docker driver
kube-start:
    minikube start --driver=docker

# Apply kubernetes manifests
kube-up:
    kubectl apply -f deployments/kubernetes/namespace.yaml
    kubectl apply -f deployments/kubernetes

# Delete kubernetes manifests
kube-down:
    kubectl delete -f deployments/kubernetes --ignore-not-found
    kubectl delete -f deployments/kubernetes/namespace.yaml --ignore-not-found

# Show kubernetes status
kube-status:
    @echo "🔍 Checking Pods in payment-gateway..."
    @kubectl get pods -n payment-gateway
    @echo -e "\n🔍 Checking Services in payment-gateway..."
    @kubectl get svc -n payment-gateway
    @echo -e "\n🔍 Checking PVCs in payment-gateway..."
    @kubectl get pvc -n payment-gateway
    @echo -e "\n🔍 Checking Jobs in payment-gateway..."
    @kubectl get jobs -n payment-gateway

# Tunnel minikube services
kube-tunnel:
    minikube tunnel

# Run auth service tests
test-auth:
    @APP_ENV=development go test service/auth/tests/... -v

# Build all Go service binaries (from api-gateway to withdraw)
build-all:
    @echo "🚀 Building all service binaries..."
    @mkdir -p bin
    @for service in {{SERVICES}}; do \
        if [ -d "service/$service" ]; then \
            echo "📦 Building $service..."; \
            (cd service/$service && go build -o ../../bin/$service $([ -f cmd/main.go ] && echo "cmd/main.go" || echo "main.go")) || exit 1; \
        fi \
    done
    @echo "✅ All binaries built in bin/ directory."

# Run go mod tidy for all services
tidy-all:
    @echo "🧹 Tidying all service modules..."
    @for service in {{SERVICES}}; do \
        if [ -d "service/$service" ]; then \
            echo "📦 Tidying $service..."; \
            (cd service/$service && go mod tidy) || echo "⚠️ Failed to tidy $service"; \
        fi \
    done
    @echo "✅ All services tidied."
