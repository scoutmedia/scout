dev:
	docker compose up gluetun -d && docker compose up scout
prod:
	env=production docker compose up --build -d
build:
	docker compose up gluetun -d && docker compose up --build scout
