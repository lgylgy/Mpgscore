VERSION=1.0
HUB=lgylgy
NAME=mpg-service
SERVICE=$(HUB)/$(NAME):$(VERSION)

build:
	docker build -t $(SERVICE) . --no-cache

run:
	docker run -it -p 3052:3052 --network host --rm $(SERVICE)

remove:
	docker rmi $(SERVICE)

push:
	docker push $(SERVICE)

deploy:
	kubectl apply -f deployment.yaml
	kubectl apply -f service.yaml