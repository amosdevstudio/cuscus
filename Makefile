.PHONY: build run bar

build:
	rm -rf build
	go build -o build/cuscus src/*.go

run:
	build/cuscus

bar: build run
