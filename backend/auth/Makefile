dev: up
	go run ./cmd/main/main.go --config=./config/config_dev.json

migrator: up
	go run ./cmd/migrator/migrator.go

create_migration:
	migrate create -ext sql -dir ./storage/migrations/ -seq 


up:
	docker compose up -d

down:
	docker compose down
