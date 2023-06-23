build:
	@go build -o ./bin/app.out .

init: build
	@./bin/app.out --initdb

run: build doc
	@./bin/app.out

test:
	@go test ./...

doc:
	@swag init -g delivery/http//main.go --output delivery/http/docs
