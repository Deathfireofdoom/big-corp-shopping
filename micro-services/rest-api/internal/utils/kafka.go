package utils

import (
    "github.com/segmentio/kafka-go"
    "log"
    "encoding/json"
)

// ToKafkaMessage handles the logic to convert a primivite type(most likely a struct)
// to a message that can be sent to the kafka-bus to a specific topic.
func ToKafkaMessage(value interface{}) (*kafka.Message, error) {
    jsonBytes, err := json.Marshal(value)
    if err != nil {
        log.Printf("could not json-marshal Interface{}: %v", err)
        return nil, err
    }

    kafkaMsg := &kafka.Message{
        Value: jsonBytes,
    }
    log.Printf("HERE: %s", jsonBytes)

    return kafkaMsg, nil
}