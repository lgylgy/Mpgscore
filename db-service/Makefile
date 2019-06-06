VERSION=1.0
HUB=lgylgy
NAME=db-service
SERVICE=$(HUB)/$(NAME):$(VERSION)

build:
	docker build -t $(SERVICE) . --no-cache

run:
	docker run -p 3051:3051 --rm $(SERVICE)

remove:
	docker rmi $(SERVICE)

push:
	docker push $(SERVICE)

deploy:
	kubectl apply -f deployment.yaml
	kubectl apply -f service.yaml