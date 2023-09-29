package main

import (
	"fmt"
	"log"
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
	scheduler.Every(15).Seconds().Do(func() {
		Scrap(kafkaClient)
	})
	scheduler.StartBlocking()
}

// Scrap scraps currencies variations
func Scrap(kc *producer.KafkaClient) error {
	urls := []string{"https://cointelegraph.com/rss", "https://cryptoslate.com/feed/", "https://www.btcethereum.com/blog/feed/"}
	var articleEvents []common.NewsEvent

	for _, v := range urls {
		articles, err := scrapper.ScrapeRSSFeed(v)
		if err != nil {
			return fmt.Errorf("could not scrap articles: %w", err)
		}
		articleEvents = append(articleEvents, articles...)
	}

	// Create and send messages with the producer
	err := kc.PushArticlesEvents(articleEvents)
	if err != nil {
		return fmt.Errorf("failed to push currency events: %w", err)
	}

	return nil
}
