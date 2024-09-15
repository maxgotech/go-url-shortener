package producer

import (
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type KafkaConfig struct {
	Username  string
	Password  string
	Bootstrap string
}

func Producer(config KafkaConfig, topic string) *kafka.Writer {
	mechanism := plain.Mechanism{
		Username: config.Username,
		Password: config.Password,
	}

	sharedTransport := &kafka.Transport{
		SASL: mechanism,
	}


	w := &kafka.Writer{
		Addr:      kafka.TCP(config.Bootstrap),
		Topic:     topic,
		Balancer:  &kafka.LeastBytes{},
		Transport: sharedTransport,
		RequiredAcks: kafka.RequireAll,
	}

	return w
}
