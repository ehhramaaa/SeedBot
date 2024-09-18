
build:
	docker build -t seedbot .

up:
	docker-compose up -d

down:
	docker-compose down

delete:
	docker rmi seedbot --force

.PHONY: build up down delete