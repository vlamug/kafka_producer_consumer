# Kafka, Producer, Consumer in Docker

### Create new network

```bash
$> docker network create app-tier --driver bridge
```

### Run ZooKeeper

```bash
$> docker run -d --name zookeeper-server \
       --network app-tier \
       -e ALLOW_ANONYMOUS_LOGIN=yes \
       bitnami/zookeeper:latest
```

### Run Kafka

```bash
$> docker run -d --name kafka-server \
       --network app-tier \
       -e ALLOW_PLAINTEXT_LISTENER=yes \
       -e KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper-server:2181 \
       bitnami/kafka:latest
```

### Run Producer

```bash
$> cd ./producer && make docker-run
```

### Run Consumer

```bash
$> cd ./consumer && make docker-run
```