# build-backend:
# 	docker buildx build . -t ubombar/daar-backend
# run-backend:
# 	docker run --rm -it -p 8080:8080 ubombar/daar-backend


start:
	docker-compose -f docker-compose.dev.yaml up

build:
	docker-compose -f docker-compose.dev.yaml build

stop:
	docker-compose -f docker-compose.dev.yaml down