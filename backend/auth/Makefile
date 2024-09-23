# Создать и запустить контейнеры со сборкой 
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
init_db: migrate users
	
# Применить миграции к базе данных
migrate: start
	docker run --rm --network auth_database -v ./:/code -w /code golang:latest go run ./cmd/migrator/migrator.go --config=./config/config_docker_dev.json

# Добавить в базу данных пользователей из csv файла
users: start
	docker run --rm --network auth_database -v ./:/code -w /code golang:latest go run ./cmd/fixture/user_data.go --config=./config/config_docker_dev.json

# Запустить тестирование проекта. Перед запуском обязательно применить команду make init_db
test: start
	.  tests/venv/bin/activate \
	&& pytest -s

# Сгенерировать описание API по комментариям в коде
swagger:
	docker run --rm -v ./:/code ghcr.io/swaggo/swag:latest init -g ./cmd/main/main.go -o ./api

