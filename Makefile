DC = docker compose -p ticket-system -f ./.docker/docker-compose.yml

.PHONY: up

up: docker-up

build: build-event
	${DC} build

down:
	${DC} down

docker-up:
	${DC} up -d --force-recreate

build-event:
	make -C event -f Makefile build