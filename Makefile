DB_PORT ?= 5432
DB_NAME ?= sns
IGNORE ?= "github.com/mdanialr/sns_backend/internal/database/sql"

clean-docker:
	docker stop postgre14-for-testing
	docker rm postgre14-for-testing

setup-docker:
	docker run --name postgre14-for-testing -e POSTGRES_PASSWORD=postgres -p "127.0.0.1:${DB_PORT}:5432" -d postgres:14-alpine

create-db:
	docker exec -it postgre14-for-testing createdb --username=postgres sns

drop-db:
	docker exec -it postgre14-for-testing dropdb --username=postgres sns

migrate:
	migrate -path internal/database/migration -database "postgresql://postgres:postgres@127.0.0.1:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migrate-down:
	migrate -path internal/database/migration -database "postgresql://postgres:postgres@127.0.0.1:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down -all

sql:
	sqlc generate

mock:
	mockgen -package mockdb -destination internal/database/mock/db.go github.com/mdanialr/sns_backend/internal/database/sql SNS

test-cover:
	go test -cover `go list ./... | grep -v ${IGNORE}`

test:
	go test -v ./...

build:
	go build -o bin/sns_backend main.go

.PHONY: clean-docker setup-docker create-db drop-db migrate migrate-down sql mock test-cover test build
