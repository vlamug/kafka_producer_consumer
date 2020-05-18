package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/alecthomas/kingpin"
)

const (
	defaultKafkaBrokers = "kafka-server:9092"
	defaultTopicName    = "test_topic_name"
)

func main() {

	var (
		app     = kingpin.New("Kafka Producer", "Produces random messages to Kafka.")
		brokers = app.Flag("brokers", "The Kafka brokers to connect to, as a comma separated list.").Default(defaultKafkaBrokers).String()
		topic   = app.Flag("topic", "The topic name to write to.").Default(defaultTopicName).String()
	)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("could not parse arguments: %s\n", err)
	}

	if *brokers == "" {
		log.Fatalln("no broker hosts set")
	}
	if *topic == "" {
		log.Fatalln("no topic name set")
	}

	producer, err := initProducers(strings.Split(*brokers, ","))
	if err != nil {
		log.Fatalf("could not init kafka producer: %s\n", err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("coudl not close producer: %s\n", err)
		}
	}()

	for {
		// run sending messages
		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: *topic,
			Value: sarama.StringEncoder(time.Now().String()),
		})
		if err != nil {
			log.Fatalf("could not send message: %s\n", err)
		}

		fmt.Printf("message was sent successfully to partition: %d and offset: %d\n", partition, offset)

		// sleep 10 seconds
		time.Sleep(time.Second * 10)
	}
}

func initProducers(brokerHosts []string) (sarama.SyncProducer, error) {
	var config = sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	return sarama.NewSyncProducer(brokerHosts, config)
}
