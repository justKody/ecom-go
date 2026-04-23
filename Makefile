.PHONY: build run

build:
	go build -o build/ecom-go cmd/main.go

run: build
	./build/ecom-go

