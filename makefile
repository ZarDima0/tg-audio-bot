# ---------------------------------------
# Variables
# ---------------------------------------
IMAGE_NAME=tg-audio-bot
CONTAINER_NAME=tg-audio-bot
BINARY_NAME=bot

# ---------------------------------------
# Запуск бота: билд и запуск Docker
# ---------------------------------------
run:
	@echo "🚀 Building Docker image and starting bot..."
	docker build -t $(IMAGE_NAME) .
	docker run -d \
		--name $(CONTAINER_NAME) \
		--env-file ./.env \
		-v ./downloads:/app/downloads \
		$(IMAGE_NAME)
	@echo "📦 Bot is running! Use 'make logs' to see logs."

# ---------------------------------------
# Просмотр логов
# ---------------------------------------
logs:
	docker logs -f $(CONTAINER_NAME)

# ---------------------------------------
# Проверка кода линтером
# ---------------------------------------
lint:
	@golangci-lint run ./...

# ---------------------------------------
# Очистка: остановка контейнера и удаление образа
# ---------------------------------------
clean:
	@echo "🧹 Stopping container and removing image..."
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true
	docker rmi $(IMAGE_NAME) || true
	rm -f $(BINARY_NAME)
	@echo "✅ Cleanup complete!"