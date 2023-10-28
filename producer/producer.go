package producer

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tobg8/crypto-viz/internal/repository"
)

func CreateProducer() (*repository.KafkaClient, error) {
	p, err := kafka.NewProducer(getKafkaConf())
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	return &repository.KafkaClient{
		Producer: p,
	}, nil
}
