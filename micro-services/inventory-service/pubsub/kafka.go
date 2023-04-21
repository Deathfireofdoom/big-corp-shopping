package pubsub


import (
	"github.com/segmentio/kafka-go"
	"encoding/json"
	"log"
	"context"
	"inventory/entity"
)

type ConsumerKafka struct {
	r *kafka.Reader
}

// TODO make it read from the last consumed.
func NewConsumerKafka(brokers []string, topic string, partition int) *ConsumerKafka {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic: topic,
		Partition: partition,
		MinBytes: 10e3, // 10KB maybe change
		MaxBytes: 10e6, // 10MB
	})

	return &ConsumerKafka{r: r}
}

func (c *ConsumerKafka) Consume(ctx context.Context, ch chan<- entity.InventoryRequest, failOnError bool) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.r.ReadMessage(ctx)
			if err != nil {
				if failOnError {
					return err
				} else {
					log.Printf("Error reading message: %v", err)
				}
			}
			
			out := entity.InventoryRequest{}
			if err := c.convertMessage(msg, &out); err != nil {
				if failOnError {
					return err
				} else {
					log.Printf("Error converting message: %v", err)
				}
			}
			ch <- out
		}
	}
}



func (c *ConsumerKafka) convertMessage(msg kafka.Message, out *entity.InventoryRequest) error {
	if err := json.Unmarshal(msg.Value, out); err != nil {
		return err
	}
	return nil
}



type ProducerKafka struct {
	w *kafka.Writer
}

func NewProducerKafka(brokers []string, topic string) *ProducerKafka {
	w := kafka.NewWriter(
		kafka.WriterConfig{
			Brokers: brokers,
			Topic: topic,
		},
	)
	return &ProducerKafka{
		w: w,
	}
}

func (p *ProducerKafka) Produce(ctx context.Context, ch <-chan entity.InventoryResult, failOnError bool) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case v := <-ch:
			// Convert to struct to kafka message
			msg, err := p.generateKafkaMsg(v)
			if err != nil {
				if failOnError{
					return err
				} else {
					log.Printf("Publish | Error on converting to kafka-msg: %v", err)
				}
			}

			err = p.w.WriteMessages(ctx, msg)
			if err != nil {
				if failOnError{
					return err
				} else {
					log.Printf("Publish | Error on publish kafka-msg: %v", err)
				}
			}

			log.Printf("DEBUG: published message.")
		}
	}
}

func (p *ProducerKafka) generateKafkaMsg(v entity.InventoryResult) (kafka.Message, error) {
	jsonMsg, err := json.Marshal(v)
	if err != nil{
		return kafka.Message{}, err 
	}

	return kafka.Message{
		Key: []byte("payload"),
		Value: jsonMsg,
	}, nil 
}
