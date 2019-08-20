run:
	go run ./cmd/main.go

up:
	docker-compose up -d

down:
	docker-compose down

dep:
	go get -d ./...