up:
	[ -d "event" ] && cd event && make up

up-debug:
	[ -d "event" ] && cd event && make up-debug

build:
	[ -d "event" ] && cd event && make build

build-debug:
	[ -d "event" ] && cd event && make build-debug

down:
	[ -d "event" ] && cd event && make down