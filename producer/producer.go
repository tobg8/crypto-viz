package producer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tobg8/crypto-viz/scrapper"
)

type KafkaClient struct {
	Producer *kafka.Producer
}

func CreateProducer() (*KafkaClient, error) {
	p, err := kafka.NewProducer(getKafkaConf())
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	return &KafkaClient{
		Producer: p,
	}, nil
}

func (kc KafkaClient) PushCurrencyEvents(currencies []scrapper.Currency) error {
	topic := "currencies_updates"
	for i := 0; i < len(currencies); i++ {
		currencyJSON, err := json.Marshal(currencies[i])
		if err != nil {
			return fmt.Errorf("failed to marshal currency to JSON: %w", err)
		}

		err = kc.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(currencies[i].Name),
			Value:          currencyJSON,
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to produce kafka message: %w", err)
		}
	}
	kc.Producer.Flush(15 * 1000)
	fmt.Printf("%v currencies events sent \n", len(currencies))
	return nil
}
