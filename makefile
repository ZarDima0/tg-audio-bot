# ---------------------------------------
# Variables
# ---------------------------------------
IMAGE_NAME=tg-audio-bot
CONTAINER_NAME=tg-audio-bot
BINARY_NAME=bot

# ---------------------------------------
# Lint
# ---------------------------------------
lint:
	golangci-lint run ./...

# ---------------------------------------
# Build
# ---------------------------------------
build:
	go build -o $(BINARY_NAME) ./cmd/bot

# ---------------------------------------
# Docker
# ---------------------------------------
docker-build:
	docker build -t $(IMAGE_NAME) .

# ---------------------------------------
# Clean
# ---------------------------------------
clean:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true
	docker rmi $(IMAGE_NAME) || true
	rm -f $(BINARY_NAME)