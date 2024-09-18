start:
	[ -d "event" ] && cd event && make build && make up && make migrate
	[ -d "booking" ] && cd booking && make build && make up && make migrate
	[ -d "checkout" ] && cd checkout && make build && make up && make migrate

start-debug:
	[ -d "event" ] && cd event && make build-debug && make up-debug && make migrate
	[ -d "booking" ] && cd booking && make build-debug && make up-debug && make migrate
	[ -d "checkout" ] && cd checkout && make build-debug && make up-debug && make migrate

stop:
	[ -d "event" ] && cd event && make down
	[ -d "booking" ] && cd booking && make down
	[ -d "checkout" ] && cd checkout && make down