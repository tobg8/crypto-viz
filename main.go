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
	usecase.InitMapPrice()
	usecase.InitMapOhlc()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// init producer
	kafkaClient, err := producer.CreateProducer()
	if err != nil {
		log.Print(err)
	}

	usecase.HandleListing(kafkaClient)
	if err != nil {
		log.Print(err)
	}

	time.Sleep(2 * time.Minute)

	// init cron jobs
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(10).Minutes().Do(func() {
		usecase.HandlePrices(kafkaClient)
		if err != nil {
			log.Print(err)
		}
		usecase.HandleOhlc(kafkaClient)
		if err != nil {
			log.Print(err)
		}
	})

	scheduler.Every(5).Minute().Do(func() {
		usecase.HandleListing(kafkaClient)
		if err != nil {
			log.Print(err)
		}
		usecase.HandleNews(kafkaClient)
		if err != nil {
			log.Print(err)
		}
	})

	scheduler.StartBlocking()
}
