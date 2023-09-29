package producer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tobg8/crypto-viz/common"
)

type KafkaClient struct {
	Producer *kafka.Producer
}

func CreateProducer() (*KafkaClient, error) {
	p, err := kafka.NewProducer(getKafkaConf())
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &KafkaClient{
		Producer: p,
	}, nil
}

func (kc KafkaClient) PushArticlesEvents(events []common.NewsEvent) error {
	topic := "news"

	if len(events) == 0 {
		log.Print("0 articles to push")
		return nil
	}

	for i := 0; i < len(events); i++ {
		articlesJSON, err := json.Marshal(events[i])
		if err != nil {
			return fmt.Errorf("failed to marshal events to JSON: %w", err)
		}

		err = kc.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(events[i].ID + " " + events[i].Title),
			Value:          articlesJSON,
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to produce kafka message: %w", err)
		}
	}
	kc.Producer.Flush(15 * 1000)
	log.Printf("%v articles events sent \n", len(events))

	return nil
}
