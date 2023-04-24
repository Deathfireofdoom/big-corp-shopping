package kafka_service

import (
	"log"
	"github.com/segmentio/kafka-go"
	"context"
	"time"
)

type KafkaCommunicator struct {
	name string
	consumer *kafka.Reader
	producer *kafka.Writer
	inChannel chan kafka.Message
	outChannel chan kafka.Message
}

func NewKafkaCommunicator(name, consumeTopic, produceTopic string, consumePartition int, inChannel, outChannel chan kafka.Message, brokers []string) (*KafkaCommunicator, error) {
	// create a producer for publishing messages to kafka
	producer, err := NewKafkaProducer(produceTopic, brokers)
	if err != nil {
		log.Printf("[error] could not create kafka-producer for Communicator %s", name)
		log.Printf("[error] %v", err)
		panic(err)
	}

	// creates a consumer for consuming messages from kafka
	consumer := NewKafkaConsumer(consumeTopic, consumePartition, brokers)
	
	return &KafkaCommunicator{
		name: name,
		consumer: consumer,
		producer: producer,
		inChannel: inChannel,
		outChannel: outChannel,
	}, nil
}

func (kc *KafkaCommunicator) Run(ctx context.Context) {
	// run functions to listen on incoming kafka messages and incoming channel messages.
	log.Printf("[debug] starting Communicator named %s", kc.name)
	go kc.consume(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Printf("[debug] recevied cancel signlar from context, shutting down Communicator %s", kc.name)
			return
		
		case msg := <- kc.inChannel:
			log.Printf("[debug] recieved inside-response on Communicator %s", kc.name)
			ctxSend, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			
			// sends message to kafka-topic
			kc.producer.WriteMessages(ctxSend, msg)
		}
	}
}


func (kc *KafkaCommunicator) consume(ctx context.Context) {
	log.Printf("[debug] consuming message from %s", kc.consumer.Config().Topic)
	for {
		msg, err := kc.consumer.ReadMessage(ctx)
		log.Printf("[debug] recieved a message on topic %s", kc.consumer.Config().Topic)
		log.Println("[debug] recieved outside-request on Communicator %s", kc.name)
			if err != nil {
				log.Println("[warning] got not ok from kafka, closing service")
				panic("got not ok from kafka, panicing..")
			}
		

		// sending kafka message for further processing. 
		kc.outChannel <- msg
	}
	log.Printf("[warning] turning of kafka-consumer")
}