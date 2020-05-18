DOCKER := docker

build:
	GOOS=linux go build -o kafka_producer_app main.go

docker-run: build
	echo "RUN"
	$(DOCKER) rm -f kafka_producer 2>/dev/null ; true
	docker build -t kafka_producer_image .
	docker run -d --name kafka_producer --network kafka-app kafka_producer_image

docker-logs:
	docker logs -f kafka_producer
