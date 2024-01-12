package usecase

import (
	"fmt"

	"github.com/tobg8/crypto-viz/internal/repository"
)

// HandleListing will currency listing from API and send them through kafka
func HandleListing(k *repository.KafkaClient) error {
	// fetch listing
	listing := repository.FetchListing()
	if len(*listing) == 0 {
		return fmt.Errorf("no new listing")
	}

	err := k.PushListing(*listing)
	if err != nil {
		return fmt.Errorf("could not send listing: %w", err)
	}
	return nil
}
