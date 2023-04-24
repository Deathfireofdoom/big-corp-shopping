package kafka_service

import (
    "time"

    "github.com/segmentio/kafka-go"
)

func NewKafkaConsumer(topic string, partition int, brokers []string) *kafka.Reader {
    // create a new Kafka reader with the specified configuration
    reader := kafka.NewReader(kafka.ReaderConfig{
        Topic:     topic,
        Partition: partition,
        MinBytes:  10e3,
        MaxBytes:  10e6,
        Brokers:   brokers,
        MaxWait:   1 * time.Second,
    })
    return reader
}