dev: up
	go run ./cmd/main/main.go --config=./config/config_dev.json

init_db: migrate users
	
migrate: 
	go run ./cmd/migrator/migrator.go
	make users

users:
	go run ./cmd/fixture/user_data.go
create_migration:
	migrate create -ext sql -dir ./storage/migrations/ -seq 

# tests; before it run migrate
test: up
	.  tests/venv/bin/activate \
	&& pytest -s

swagger:
	docker run --rm -v ./:/code ghcr.io/swaggo/swag:latest init -g ./cmd/main/main.go -o ./api

up:
	docker compose up -d

down:
	docker compose down
