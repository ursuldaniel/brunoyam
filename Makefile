build:
	@go build -o bin/app ./cmd/rest-api-practice/main.go

run: build
	@./bin/app