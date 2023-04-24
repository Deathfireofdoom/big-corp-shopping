package config

import (
	"fmt"
	"github.com/spf13/viper"
)


var KafkaInventoryRequest kafkaInventoryRequest
type  kafkaInventoryRequest struct {
	RequestTopic 	string 		`mapstructure:"request_topic"`
	ResponseTopic 	string 		`mapstructure:"response_topic"`
	Brokers 		[]string	`mapstructure:"brokers"`
}

var KafkaCartRequest kafkaCartRequest
type  kafkaCartRequest struct {
	RequestTopic 	string 		`mapstructure:"request_topic"`
	ResponseTopic 	string 		`mapstructure:"response_topic"`
	Brokers 		[]string	`mapstructure:"brokers"`
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

	err = viper.UnmarshalKey("kafka_inventory_request", &KafkaInventoryRequest)
	if err != nil {
		panic(fmt.Errorf("could not unmarshal key kafka_inventory_request: %s", err))
	}

	err = viper.UnmarshalKey("kafka_cart_request", &KafkaCartRequest)
	if err != nil {
		panic(fmt.Errorf("could not unmarshal key kafka_cart_request: %s", err))
	}

}