
start:
	docker-compose -f docker-compose.dev.yaml up

build:
	docker-compose -f docker-compose.dev.yaml build

remove:
	docker-compose -f docker-compose.dev.yaml down