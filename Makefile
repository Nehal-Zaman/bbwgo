build:
	@go build -o bin/bbwgo -ldflags="-s -w"

run: build
	@./bin/bbwgo