DC=docker compose -p ticket-system --env-file ./.env -f ./.docker/docker-compose.yml

.DEFAULT_GOAL := help

### docker

start: start-shared-containers ## build and run containers
	[ -d "event" ] && cd event && make build && make up && make migrate
	[ -d "booking" ] && cd booking && make build && make up && make migrate
	[ -d "checkout" ] && cd checkout && make build && make up && make migrate

start-debug: start-shared-containers ## build and run containers with debugger
	[ -d "event" ] && cd event && make build-debug && make up-debug && make migrate
	[ -d "booking" ] && cd booking && make build-debug && make up-debug && make migrate
	[ -d "checkout" ] && cd checkout && make build-debug && make up-debug && make migrate

stop:
	[ -d "event" ] && cd event && make down
	[ -d "booking" ] && cd booking && make down
	[ -d "checkout" ] && cd checkout && make down
	stop-shared-containers


### helpers

start-shared-containers: env
	${DC} up -d

stop-shared-containers:
	${DC} down

help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[32m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

env:
	if [ ! -f .env ]; then \
		cp .env.dist .env; \
	fi