package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tobg8/crypto-viz/common"
)

func TestIsListingNew(t *testing.T) {
	t.Run("when listing is not new", func(t *testing.T) {
		in := []common.ListingEvent{
			{
				CurrentPrice: 1292.921,
				Symbol:       "hello",
			},
			{
				CurrentPrice: 129992,
				Symbol:       "tata",
			},
			{
				CurrentPrice: 122132,
				Symbol:       "tate",
			},
		}

		price := 1292.921

		expect := false
		out := isListingNew(in, &price)
		assert.Equal(t, expect, out)
	})
	t.Run("when listing is new", func(t *testing.T) {
		in := []common.ListingEvent{
			{
				CurrentPrice: 1292.921,
				Symbol:       "hello",
			},
			{
				CurrentPrice: 129992,
				Symbol:       "tata",
			},
			{
				CurrentPrice: 122132,
				Symbol:       "tate",
			},
		}

		price := 1292.1

		expect := true
		out := isListingNew(in, &price)
		assert.Equal(t, expect, out)
	})
}
