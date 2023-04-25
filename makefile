APP_NAME=app
DOCKER_TAG=latest

.PHONY: build docker clean run-mongo

build:
	go build -o app .

run:
	go run cmd/main.go

docker:
	docker build -t app:latest .

clean:
	rm -f app

run-mongo:
	docker run -p 27017:27017 -d --rm mongo:latest
