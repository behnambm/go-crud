build:
	@go build -o ./bin/app.out .

init: build
	@./bin/app.out --initdb

run: build
	@./bin/app.out

