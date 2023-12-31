dev:
	docker compose up gluetun -d && docker compose up scout
prod:
	env=production docker compose up --build -d
build:
	docker compose up --build scout
run:
	env=development go run main.go