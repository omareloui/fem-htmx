dev:
	@air

run:
	@go run cmd/fem-htmx/main.go

build: 
	@go build -o .output/fem-htmx cmd/fem-htmx/main.go

start: 
	@./.output/fem-htmx
