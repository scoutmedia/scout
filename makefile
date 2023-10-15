run:
	go run main.go
dev:
	docker compose up --build -d
prod:
	env=production docker compose up --build -d
build:
	go build && sudo ./scout
