package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/tobg8/crypto-viz/kafkaclient"
	"github.com/tobg8/crypto-viz/scrapper"
)

func main() {
	Init()
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// TODO CONSUMER
	// init Consumer

	// init producer
	kafkaClient, err := kafkaclient.CreateProducer()
	if err != nil {
		log.Print(err)
	}

	// init cron jobs
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(30).Seconds().Do(func() {
		Scrap(kafkaClient)
	})
	scheduler.StartBlocking()
}

// Scrap scraps currencies variations about top 100 crypto currencies.
func Scrap(kc *kafkaclient.KafkaClient) error {
	currencies, err := scrapper.ScrapCurrencies("https://www.coingecko.com/fr")
	if err != nil {
		log.Print(err)
	}
	// create and send message with producer
	err = kc.PushCurrencyEvents(currencies)
	if err != nil {
		return fmt.Errorf("failed to push currency events: %w", err)
	}

	// newCurrencies, err := scrapper.ScrapNewCurrencies("https://www.coingecko.com/fr/new-cryptocurrencies")
	// if err != nil {
	// 	log.Print(err)
	// }
	// log.Printf("NEW CURRENCIES: %v", newCurrencies)
	return nil
}
