package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/alecthomas/kingpin"
)

const (
	defaultKafkaBrokers            = "kafka-server:9092"
	defaultTopicName               = "test_topic_name"
	defaultConsumerGroupDefinition = "test_consumer_group_definition"
)

func main() {

	var (
		app     = kingpin.New("Kafka Producer", "Consumes messages from Kafka.")
		brokers = app.Flag("brokers", "The Kafka brokers to connect to, as a comma separated list.").Default(defaultKafkaBrokers).String()
		topics  = app.Flag("topics", "The topic name to read from.").Default(defaultTopicName).String()
	)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("could not parse arguments: %s\n", err)
	}

	if *brokers == "" {
		log.Fatalln("no broker hosts set")
	}
	if *topics == "" {
		log.Fatalln("no topic name set")
	}

	consumer, err := initConsumer(strings.Split(*brokers, ","))
	if err != nil {
		log.Fatalf("could not init kafka producer: %s\n", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("coudl not close producer: %s\n", err)
		}
	}()

	for {
		// run sending messages
		err := consumer.Consume(context.Background(), strings.Split(*topics, ","), &ConsumerHandler{})
		if err != nil {
			log.Fatalf("could not send message: %s\n", err)
		}
	}
}

type ConsumerHandler struct{}

func (c *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}

	return nil
}

func initConsumer(brokerHosts []string) (sarama.ConsumerGroup, error) {
	var config = sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return sarama.NewConsumerGroup(brokerHosts, defaultConsumerGroupDefinition, config)
}
