package kafka_service

import (
	"github.com/segmentio/kafka-go"
)

// NewKafkaProduct handles the logic to create a new producer that can
// be used to send message to kafka-bus.
func NewKafkaProducer(topic string, brokerURLs[]string) (*kafka.Writer, error) {
    // Configuring 
    config := kafka.WriterConfig{
        Brokers: brokerURLs,
        Topic:   topic,
    }

    // create a new Kafka producer instance
    producer := kafka.NewWriter(config)
    return producer, nil
}

