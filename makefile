# ---------------------------------------
# Variables
# ---------------------------------------
IMAGE_NAME=tg-audio-bot
CONTAINER_NAME=tg-audio-bot

# ---------------------------------------
# Команда 1: Собрать и поднять приложение
# ---------------------------------------
up:
	docker-compose up -d --build

# ---------------------------------------
# Команда 2: Остановить и почистить ТОЛЬКО ПРОЕКТ
# ---------------------------------------
down:
	docker-compose down --rmi all --volumes --remove-orphans

# ---------------------------------------
# Вспомогательные команды
# ---------------------------------------
logs:
	docker-compose logs -f

stop:
	docker-compose stop

status:
	docker-compose ps

restart: stop up

# ---------------------------------------
# Локальная разработка
# ---------------------------------------
build:
	go build -o bot ./cmd/bot

run: build
	./bot

lint:
	golangci-lint run ./...