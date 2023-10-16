run:
	go run main.go
dev:
	docker compose up --build
prod:
	env=production docker compose up --build -d
build:
	go build && sudo ./scout
