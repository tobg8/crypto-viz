package usecase

import (
	"fmt"
	"log"

	"github.com/tobg8/crypto-viz/common"
	"github.com/tobg8/crypto-viz/internal/repository"
)

var lastBTCValue float64

// HandleListing will currency listing from API and send them through kafka
func HandleListing(k *repository.KafkaClient) error {
	// fetch listing
	listing := repository.FetchListing()
	if len(*listing) == 0 {
		return fmt.Errorf("no new listing")
	}
	newListing := isListingNew(*listing, &lastBTCValue)
	if !newListing {
		log.Print("no new listing to send")
		return nil
	}

	lastBTCValue = (*listing)[0].CurrentPrice
	err := k.PushListing(*listing)
	if err != nil {
		return fmt.Errorf("could not send listing: %w", err)
	}
	return nil
}

// isListingNew returns whether the listing has already been processed
func isListingNew(l []common.Listing, price *float64) bool {
	return l[0].CurrentPrice != *price
}
