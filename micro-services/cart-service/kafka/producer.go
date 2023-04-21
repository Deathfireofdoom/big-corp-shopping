package kafka

import (
	"github.com/segmentio/kafka-go"
	"encoding/json"
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

// ToKafkaMessage handles the logic to convert a primivite type(most likely a struct)
// to a message that can be sent to the kafka-bus to a specific topic.
func ToKafkaMessage(topic string, value interface{}) (*kafka.Message, error) {
    jsonBytes, err := json.Marshal(value)
    if err != nil {
        return nil, err
    }

    kafkaMsg := &kafka.Message{
        Topic: topic,
        Value: jsonBytes,
    }

    return kafkaMsg, nil
}