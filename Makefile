BASEDIR := ${CURDIR}
GO      := $(shell which go)
GOPATH  := $(shell go env GOPATH)
GOBIN   := $(GOPATH)/bin

IMAGE_NAME = kaze-image
DOCKER_BUILD_CONTEXT=.
DOCKER_FILE=Dockerfile
COMPOSE_FILE=docker-compose.yml

init:
	go mod download
	go mod tidy

run:
	go run ./cmd/kaze/

test/unit:
	go test -race ./...

test/vet:
	$(GO) vet ./...

clean:
	rm -rf build || true
	mkdir build

build: clean
	go build -o build/kaze ./cmd/kaze

# Build the Docker image
docker/build:
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME) -f $(DOCKER_FILE) $(DOCKER_BUILD_CONTEXT)

# Run the Docker container
docker/run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(IMAGE_NAME)

# Start Docker Compose services
compose/build:
	@echo "Starting Docker Compose services..."
	docker-compose -f $(COMPOSE_FILE) build 

# Start Docker Compose services
compose/up:
	@echo "Starting Docker Compose services..."
	docker-compose -f $(COMPOSE_FILE) up 

# Stop Docker Compose services
compose/down:
	@echo "Stopping Docker Compose services..."
	docker-compose -f $(COMPOSE_FILE) down