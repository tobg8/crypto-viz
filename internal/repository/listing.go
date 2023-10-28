package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tobg8/crypto-viz/common"
)

func (kc KafkaClient) PushListing(a []common.Listing) error {
	topic := "listing"
	for i := 0; i < len(a); i++ {
		article, err := json.Marshal(a[i])
		if err != nil {
			return fmt.Errorf("failed to marshal articles to JSON: %w", err)
		}

		err = kc.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(a[i].Name),
			Value:          article,
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to produce kafka message: %w", err)
		}
	}

	kc.Producer.Flush(15 * 1000)
	log.Printf("%v currencies in listing sent \n", len(a))
	return nil
}

func FetchListing() *[]common.Listing {
	client := &http.Client{}
	apiKey := os.Getenv("COINGECKO_KEY")
	baseURL := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=eur&order=market_cap_desc&per_page=250&locale=fr"
	req, err := http.NewRequest("GET", baseURL, nil)
	req.Header.Set("X-CG-Demo-API-Key", apiKey)

	var response []common.Listing
	if err != nil {
		log.Printf("could not create request: %v", err)
		return &response
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("could not get URL: %v", baseURL)
		return &response
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
		return &response
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	return &response
}
