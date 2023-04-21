package kafka_service

import (
	"github.com/segmentio/kafka-go"
	"context"
	"log"
	"time"
)

// RequestCommunicator is responsible for communicate over request-response. 
// Reading request from Kafka topic and then publishing response to another kafka topic.
type KafkaCommunicator struct {
	name string 
	outsideConsumer *kafka.Reader
	outsideProducer *kafka.Writer
	inChannel chan kafka.Message
	outChannel chan kafka.Message
}


func NewKafkaCommunicator(name, incommingTopic, outgoingTopic string, inCommingPartition int, inChannel, outChannel chan kafka.Message) (*KafkaCommunicator, error) {
	// creates producer
	producer, err := NewKafkaProducer(outgoingTopic, config.KafkaBridge.Brokers)
	if err != nil {
		log.Printf("[error] could not create kafka-producer for Communicator %s", name)
		log.Printf("[error] %v", err)
		panic(err)
	}

	// creates consumer
	consumer := NewKafkaConsumer(incommingTopic, inCommingPartition, config.KafkaBridge.Brokers)
	return &KafkaCommunicator{
		name: name,
		outsideConsumer: consumer,
		outsideProducer: producer,
		inChannel: inChannel,
		outChannel: outChannel,
	}
}

func (kc *KafkaCommunicator) Run(ctx context.Context) {
	// run functions to listen on incoming kafka messages and incoming channel messages.
	log.Printf("[debug] starting Communicator named %s", kc.name)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[debug] recevied cancel signlar from context, shutting down Communicator %s", tc.name)
			return
		
		case msg, ok := kc.outsideConsumer.Messages():
			log.Println("[debug] recieved outside-request on Communicator %s", kc.name)
			if !ok {
				log.Println("[warning] got not ok from kafka, closing service")
				return 
			}
			
			// sending kafka message for further processing. 
			kc.outChannel <- msg.Value
		
		case msg := <- kc.outChannel:
			log.Println("[debug] recieved inside-response on Communicator %s", kc.name)
			ctxSend, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			
			// sends message to kafka-topic
			kc.outsideProducer.WriteMessages(ctxSend, *msg)
		}
	}
}
