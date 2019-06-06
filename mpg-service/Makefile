VERSION=1.0
HUB=lgylgy
NAME=mpg-service
SERVICE=$(HUB)/$(NAME):$(VERSION)

build:
	docker build -t $(SERVICE) . --no-cache

run:
	docker run -p 3052:3052 --rm $(SERVICE)

remove:
	docker rmi $(SERVICE)

push:
	docker push $(SERVICE)

deploy:
	kubectl apply -f deployment.yaml
	kubectl apply -f service.yaml