DC = docker compose -p ticket-system -f ./.docker/docker-compose.yml

.PHONY: all

up: docker-up

build:
	${DC} build

down:
	${DC} down

docker-up:
	${DC} up -d --force-recreate