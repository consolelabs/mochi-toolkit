package kafka

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (k *Kafka) RunConsumer(groupID, topic string, handler func(payload []byte) error) error {
	instanceId, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %s", err)
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               k.brokers,
		"group.id":                        groupID,
		"auto.offset.reset":               "earliest",
		"enable.partition.eof":            true,
		"go.application.rebalance.enable": true,
		"go.events.channel.enable":        true,
		"session.timeout.ms":              5 * 60 * 1000,
		"group.instance.id":               instanceId,
		"enable.auto.commit":              "false",
		"max.partition.fetch.bytes":       1048576 * 2,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %s", err)
	}

	if err := c.Subscribe(topic, nil); err != nil {
		return fmt.Errorf("failed to subscribe to topic: %s", err)
	}

	run := true
	for run {
		select {
		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				k.logger.WithField("partitions", e.String()).Info("assigned partitions")
				c.Assign(e.Partitions)

			case kafka.RevokedPartitions:
				k.logger.WithField("partitions", e.String()).Info("revoked partitions")
				c.Unassign()

			case *kafka.Message:
				if err := k.processMessage(e, handler); err != nil {
					k.logger.WithError(err).Error("failed to process message")
				}

				i := 0
				for {
					if _, err := c.CommitMessage(e); err != nil {
						if i < 10 {
							i++
							continue
						}
						return fmt.Errorf("failed to commit offset: %s", err)
					}
					break
				}

			case kafka.PartitionEOF:
				k.logger.Infof("reached %v\n", e)

			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				k.logger.WithError(e).Error("consumer error")
			}
		}
	}

	if err := c.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %s", err)
	}

	return nil
}

func (k *Kafka) processMessage(msg *kafka.Message, handler func(payload []byte) error) error {
	if err := handler(msg.Value); err != nil {
		return fmt.Errorf("failed to handle message: %s message: %s", err, string(msg.Value))
	}

	return nil
}
