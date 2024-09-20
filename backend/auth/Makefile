# dev: up
# 	go run ./cmd/main/main.go --config=./config/config_dev.json

init_db: migrate users
	
migrate: 
	docker run --rm -v .:/code -w /app golang:latest go run ./cmd/migrator/migrator.go

users:
	docker run --rm -v .:/code -w /app golang:latest go run ./cmd/fixture/user_data.go

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
