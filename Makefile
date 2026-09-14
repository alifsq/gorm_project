include .env
export
.PHONY: run build gen

go-test-user:
	go test -v ./internal/modules/user/...

run: 
	go run cmd/api/main.go

gen: 
	go run cmd/gen/main.go

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-create:
	migrate create -ext sql -dir migrations $(name)