package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tobg8/crypto-viz/common"
)

func (kc KafkaClient) FetchOhlc(currency string, rangeCurrency string) *common.OhlcResponseAPI {
	var tempResp [][]json.RawMessage
	var response common.OhlcResponseAPI

	baseURL := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%v/ohlc?vs_currency=eur&days=%v", currency, rangeCurrency)
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

	// Process ohlc
	response.Ohlc, err = processOhlcMessages(tempResp)
	if err != nil {
		log.Fatal(err)
	}

	return &response
}

func processOhlcMessages(rawMessages [][]json.RawMessage) ([]common.OhlcUnit, error) {
	var result []common.OhlcUnit
	for _, rawMsg := range rawMessages {
		var values []float64
		if len(rawMsg) != 5 {
			return nil, fmt.Errorf("invalid number of values in raw message: %v", len(rawMsg))
		}

		for _, rawValue := range rawMsg {
			var value float64
			if err := json.Unmarshal(rawValue, &value); err != nil {
				return nil, err
			}
			values = append(values, value)
		}

		timestamp := int64(values[0])
		open := values[1]
		high := values[2]
		low := values[3]
		close := values[4]

		result = append(result, common.OhlcUnit{Timestamp: timestamp, Open: open, High: high, Low: low, Close: close})
	}
	return result, nil
}

func (kc KafkaClient) PushOhlc(p common.OhlcResponseTest, r string, c string) error {
	topic := "ohlc"

	ohlc, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal ohlc to JSON: %w", err)
	}

	err = kc.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte("ohlc:" + p.Currency),
		Value:          ohlc,
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to produce kafka message: %w", err)
	}

	kc.Producer.Flush(15 * 1000)
	log.Printf("%v ohlc  sent on range %v and crypto %v \n", len(p.Ohlc), r, c)
	return nil
}
