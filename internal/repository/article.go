package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tobg8/crypto-viz/common"
)

type KafkaClient struct {
	Producer *kafka.Producer
}

func (kc KafkaClient) PushArticles(a []common.ArticleEvent) error {
	topic := "news"
	for i := 0; i < len(a); i++ {
		article, err := json.Marshal(a[i])
		if err != nil {
			return fmt.Errorf("failed to marshal articles to JSON: %w", err)
		}

		err = kc.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(a[i].Title),
			Value:          article,
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to produce kafka message: %w", err)
		}
	}

	kc.Producer.Flush(15 * 1000)
	log.Printf("%v article events sent \n", len(a))
	return nil
}

func FetchArticles() *[]common.ArticleAPI {
	var response common.APIResponse

	baseURL := "https://cryptopanic.com/api/v1/posts/?auth_token=c1b068d015189cba73a935bb42c06128b4c3e5f6"
	resp, err := http.Get(baseURL)
	if err != nil {
		log.Printf("could not get url: %v", baseURL)
		return &response.Results
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
		return &response.Results
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	return &response.Results
}
