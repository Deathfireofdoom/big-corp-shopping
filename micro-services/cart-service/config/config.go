package config

import (
	"fmt"
	"github.com/spf13/viper"
)


var KafkaConsumerConfig kafkaConsumerConfig
type kafkaConsumerConfig struct {
	Topic string
	Partition int
	Brokers []string
}

var KafkaProducerConfig kafkaProducerConfig
type kafkaProducerConfig struct {
	Brokers []string
	Topic string 
}

// LoadConfig loads the config and saves it into variables that is exported from the
// config package.
func init() {
	// Read config file with viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("could not read configuration file: %s", err))
	}

	err = viper.UnmarshalKey("kafka_consumer", &KafkaConsumerConfig)
	if err != nil {
		panic(fmt.Errorf("could not unmarshal key kafka_consumer: %s", err))
	}

	err = viper.UnmarshalKey("kafka_producer", &KafkaProducerConfig)
	if err != nil {
		panic(fmt.Errorf("could not unmarshal key kafka_producer: %s", err))
	}

}