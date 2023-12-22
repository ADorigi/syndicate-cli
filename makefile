build:
	@go build -v ./...

run:
	@./syndicate-cli

brun:
	@go build -v ./...
	@./syndicate-cli