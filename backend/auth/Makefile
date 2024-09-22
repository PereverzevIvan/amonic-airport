# dev: up
# 	go run ./cmd/main/main.go --config=./config/config_dev.json

init_db: migrate users
	
migrate: 
	docker run --rm --network auth_database -v ./:/code -w /code golang:latest go run ./cmd/migrator/migrator.go --config=./config/config_docker_dev.json

users:
	docker run --rm --network auth_database -v ./:/code -w /code golang:latest go run ./cmd/fixture/user_data.go --config=./config/config_docker_dev.json

# tests; before it run migrate
test: up
	.  tests/venv/bin/activate \
	&& pytest -s

swagger:
	docker run --rm -v ./:/code ghcr.io/swaggo/swag:latest init -g ./cmd/main/main.go -o ./api

up:
	docker compose up -d --build

down:
	docker compose down

start:
	docker compose start

stop:
	docker compose stop
