build:
	@go build -o ./bin/app.out .

init: build
	@./bin/app.out --initdb

run: doc build
	@./bin/app.out

test:
	@go test ./...

doc:
	@swag init -g delivery/http/main.go --output delivery/http/docs

docker-build: build
	@docker build -t b/book:latest .

up: docker-build
	@docker compose up

