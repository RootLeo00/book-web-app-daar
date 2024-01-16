
start:
	docker-compose -f docker-compose.dev.yaml up

build:
	docker-compose -f docker-compose.dev.yaml build

stop:
	docker-compose -f docker-compose.dev.yaml down