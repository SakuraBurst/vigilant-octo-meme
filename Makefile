.PHONY: dc run test lint

dc:
	docker-compose up --remove-orphans --build

run:
	go build -o app cmd/photos-backend/main.go && ./app


test:
	go test -race ./...

lint:
	golangci-lint run
