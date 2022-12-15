postgres:
	docker run --name postgres1 -p 5432:5432 -e POSTGRES_USER=stunna -e POSTGRES_PASSWORD=ahmed143 -d postgres:15-alpine
createdb:
	docker exec -it postgres1 createdb --username=stunna --owner=stunna simple_bank
dropdb:
	docker exec -it postgres1 dropdb
migrateup:
	migrate -path db/migration/ -database "postgresql://stunna:ahmed143@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration/ -database "postgresql://stunna:ahmed143@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration/ -database "postgresql://stunna:ahmed143@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration/ -database "postgresql://stunna:ahmed143@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

-PHONY: postgres createdb dropdb migrateup migratedown sqlc test server