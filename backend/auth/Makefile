dev: up
	go run ./cmd/main/main.go --config=./config/config_dev.json

migrate: up
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

up:
	docker compose up -d

down:
	docker compose down
