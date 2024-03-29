package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tobg8/crypto-viz/common"
)

func (kc KafkaClient) FetchPrices(currency string, rangeCurrency string) *common.PriceResponseAPI {
	var tempResp common.PriceAPI
	var response common.PriceResponseAPI

	baseURL := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%v/market_chart?vs_currency=eur&days=%v", currency, rangeCurrency)
	resp, err := http.Get(baseURL)
	if err != nil {
		log.Printf("could not get url: %v", baseURL)
		return &response
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
		return &response
	}

	err = json.NewDecoder(resp.Body).Decode(&tempResp)
	if err != nil {
		log.Fatal(err)
	}

	// Process Prices
	response.Prices, err = processRawMessages(tempResp.Prices)
	if err != nil {
		log.Fatal(err)
	}

	return &response
}

func processRawMessages(rawMessages []json.RawMessage) ([]common.PriceUnitAPI, error) {
	var result []common.PriceUnitAPI
	for _, rawMsg := range rawMessages {
		var values []interface{}
		if err := json.Unmarshal(rawMsg, &values); err != nil {
			return nil, err
		}

		// API so fucked up they have value without timestamp
		// Or timestamp without value sometimes
		if len(values) == 2 {
			timestamp, ok := values[0].(float64)
			if !ok {
				return nil, fmt.Errorf("timestamp is not a float64")
			}

			value, ok := values[1].(float64)
			if !ok {
				return nil, fmt.Errorf("value is not a float64")
			}

			result = append(result, common.PriceUnitAPI{Timestamp: int64(timestamp), Value: value})
		} else {
			return nil, fmt.Errorf("invalid number of values in raw message: %v", len(values))
		}
	}
	return result, nil
}

func (kc KafkaClient) PushPrices(p common.PriceEventTest, r string, c string) error {
	topic := "prices"

	price, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal articles to JSON: %w", err)
	}

	err = kc.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte("price:" + p.Currency),
		Value:          price,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to produce kafka message: %w", err)
	}

	kc.Producer.Flush(15 * 1000)
	log.Printf("%v prices in listing sent on range %v and crypto %v \n", len(p.Prices), r, c)
	return nil
}
