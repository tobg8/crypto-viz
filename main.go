package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/tobg8/crypto-viz/internal/usecase"
	"github.com/tobg8/crypto-viz/producer"
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

	// isListingDone := false
	// init cron jobs
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Minutes().Do(func() {
		log.Print("here we go")
		// isListingDone = true
		// if isListingDone {
		// 	usecase.HandleNews(kafkaClient)
		err := usecase.HandlePrices(kafkaClient)
		if err != nil {
			log.Print(err)
		}
		// }
		// usecase.HandleListing(kafkaClient)

	})
	scheduler.StartBlocking()
}
