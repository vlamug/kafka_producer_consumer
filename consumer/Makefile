DOCKER := docker

build:
	GOOS=linux go build -o kafka_consumer_app main.go

docker-run: build
	echo "RUN"
	$(DOCKER) rm -f kafka_consumer 2>/dev/null ; true
	docker build -t kafka_consumer_image .
	docker run -d --name kafka_consumer --network kafka-app kafka_consumer_image

docker-logs:
	docker logs -f kafka_consumer
