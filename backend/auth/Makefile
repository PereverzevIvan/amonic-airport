# Создать и запустить контейнеры со сборкой 
dev:
	docker compose up mysql phpmyadmin swagger-ui -d
	go run ./cmd/main/main.go --config=./config/config_dev.json

up:
	docker compose up -d --build

# Удалить все контейнеры
down:
	docker compose down

# Запустить контейнеры
start:
	docker compose start

# Остановить все контейнеры
stop:
	docker compose stop

# Перезапустить все контейнеры
restart: 
	docker compose restart

# Заполнить базу данных тестовой информацией (нужно запускать только один раз, после создания контейнеров)
local-init-db: local-migrate local-users
init-db: migrate users

# Применить миграции к базе данных
local-migrate:
	go run ./cmd/migrator/migrator.go --config=./config/config_dev.json
migrate: start
	docker run --rm --network auth_database -v ./:/code -w /code golang:latest go run ./cmd/migrator/migrator.go --config=./config/config_docker_dev.json

# Добавить в базу данных пользователей из csv файла
local-users:
	go run ./cmd/fixture/user_data.go --config=./config/config_dev.json
users: start
	docker run --rm --network auth_database -v ./:/code -w /code golang:latest go run ./cmd/fixture/user_data.go --config=./config/config_docker_dev.json

# Запустить тестирование проекта. Перед запуском обязательно применить команду make init_db
test: start
	.  tests/venv/bin/activate \
	&& pytest -s

# Сгенерировать описание API по комментариям в коде
swagger:
	docker run --rm -v ./:/code ghcr.io/swaggo/swag:latest init -g ./cmd/main/main.go -o ./api

