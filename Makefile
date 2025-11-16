build:
	docker compose build

up:
	docker compose up

run: build
	docker compose up

down:
	docker compose down