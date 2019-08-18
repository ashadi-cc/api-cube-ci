run:
	go get -d ./...
	go run .

up:
	docker-compose up

down:
	docker-compose down