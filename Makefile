run:
	go run main.go

build:
	go build -o task-api

docker:
	docker build -t task-api .

docker-run:
	docker run --rm -p 8081:8081 task-api
