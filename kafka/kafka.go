package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

type Kafka struct {
	brokers  string
	producer *producer
	logger   *logrus.Entry
}

type producer struct {
	producer *kafka.Producer
	ready    bool
	termChan chan bool
}

func New(brokers string, logger *logrus.Entry) *Kafka {
	if brokers == "" {
		return nil
	}
	return &Kafka{
		brokers:  brokers,
		producer: &producer{},
		logger: logrus.WithFields(logrus.Fields{
			"component": "kafka",
		}),
	}
}
