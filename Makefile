run:
	go run ./cmd/main.go

up:
	docker-compose up

down:
	docker-compose down

dep:
	go get -d ./...