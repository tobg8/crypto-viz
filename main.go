package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/tobg8/crypto-viz/common"
	"github.com/tobg8/crypto-viz/producer"
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

	// init producer
	kafkaClient, err := producer.CreateProducer()
	if err != nil {
		log.Print(err)
	}

	// init cron jobs
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(10).Seconds().Do(func() {
		Scrap(kafkaClient)
	})
	scheduler.StartBlocking()
}

// Scrap scraps currencies variations
func Scrap(kc *producer.KafkaClient) error {
	// Use a wait group to handle concurrent data fetching
	var wg sync.WaitGroup

	// Use channels to receive the results
	currenciesCh := make(chan []common.CurrencyEvent)
	// newCurrenciesCh := make(chan []common.CurrencyEvent)

	fetchCurrencies := func(url string, destCh chan<- []common.CurrencyEvent) {
		defer wg.Done()
		curr, err := scrapper.ScrapCurrencies(url)
		if err != nil {
			log.Printf("Error fetching currencies from %s: %v", url, err)
			destCh <- nil
			return
		}
		destCh <- curr
	}

	// fetchNewCurrencies := func(url string, destCh chan<- []common.CurrencyEvent) {
	// 	defer wg.Done()
	// 	newCurr, err := scrapper.ScrapNewCurrencies(url)
	// 	if err != nil {
	// 		log.Printf("Error fetching new currencies from %s: %v", url, err)
	// 		destCh <- nil
	// 		return
	// 	}
	// 	destCh <- newCurr
	// }

	// Fetch currencies and newCurrencies
	wg.Add(2)
	go fetchCurrencies("https://coinmarketcap.com/", currenciesCh)
	// go fetchNewCurrencies("https://www.coingecko.com/fr/new-cryptocurrencies", newCurrenciesCh)

	go func() {
		wg.Wait()
		close(currenciesCh)
		// close(newCurrenciesCh)
	}()

	currencies := <-currenciesCh
	// newCurrencies := <-newCurrenciesCh

	// Create and send messages with the producer
	if err := kc.PushCurrencyEvents(currencies); err != nil {
		return fmt.Errorf("failed to push currency events: %w", err)
	}

	return nil
}
