package config

import (
	"fmt"
	"github.com/spf13/viper"
)


var KafkaConfig kafkaConfig
type kafkaConfig struct {
	RequestTopic 	string 		`mapstructure:"request_topic"`
	ResponseTopic 	string 		`mapstructure:"response_topic"`
	Brokers 		[]string	`mapstructure:"brokers"`
}

// init loads the config and saves it into variables that is exported from the
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

	err = viper.UnmarshalKey("kafka", &KafkaConfig)
	if err != nil {
		panic(fmt.Errorf("could not unmarshal key kafka: %s", err))
	}
}