createdb:
	createdb --username=wahid --owner=wahid shopping

dropdb:
	dropdb --username=wahid shopping

migrateup:
	migrate -path db/migration -database "postgresql://wahid:secret@localhost:5432/shopping?sslmode=disable" -verbose up
	
migratedown:
	migrate -path db/migration -database "postgresql://wahid:secret@localhost:5432/shopping?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://wahid:secret@localhost:5432/shopping?sslmode=disable" -verbose up 1
	
migratedown1:
	migrate -path db/migration -database "postgresql://wahid:secret@localhost:5432/shopping?sslmode=disable" -verbose down 1


sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server