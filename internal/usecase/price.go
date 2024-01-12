package usecase

import (
	"fmt"
	"log"
	"sort"

	"github.com/tobg8/crypto-viz/common"
	"github.com/tobg8/crypto-viz/internal/repository"
)

var Maped map[string]int64

func InitMapPrice() {
	// Initialize the map before using it
	Maped = make(map[string]int64)
}

func HandlePrices(k *repository.KafkaClient) error {
	log.Print("new prices")
	// Je veux appeler ma func repo qui récup les prix en passant range et currency
	currencies := repository.FetchListing()
	// Je les sauvegarde
	for _, v := range *currencies {
		if _, ok := Maped[v.ID+"1D"]; !ok {
			Maped[v.ID+"1D"] = 0
		}
	}

	for i, v := range *currencies {
		// Pour l'instant on gère 10 currencies en One day
		// Parce que on prends de temps que les 5 mins de la routine
		if i == 11 {
			return nil
		}
		_, err := handle1DPrice(k, v.ID)
		if err != nil {
			log.Print(err)
			return nil
		}
	}

	return nil
}

func transformToEvent(prices *common.PriceResponseAPI, r string, currency string) common.PriceEventTest {
	var response common.PriceEventTest
	for _, v := range prices.Prices {
		if len(prices.Prices) > 0 {
			temp := common.PriceEvent{
				PriceUnitAPI: v,
				Type:         "prices",
				Range:        r,
				Currency:     currency,
			}
			response.Prices = append(response.Prices, temp)
		}
	}
	response.Currency = currency
	return response
}

func handle1DPrice(k *repository.KafkaClient, currency string) (*common.PriceResponseAPI, error) {
	prices := k.FetchPrices(currency, "1")
	if prices == nil {
		return nil, fmt.Errorf("error when fetching prices on 1D")
	}
	// Je veux trier mes arrays par timestamp,
	sortPricesByTimeStamp(prices)

	// Check if the stored timestamp for the currency is equal to or older than the latest timestamp
	if storedTimestamp, ok := Maped[currency+"1D"]; ok && storedTimestamp >= prices.Prices[0].Timestamp {
		return nil, fmt.Errorf("no new prices to send on 1D")
	}

	// J'enlève ceux qui sont plus grands
	removeOlderValues(Maped[currency+"1D"], prices)

	// Je set le lastTimeStamp avec le plus récent
	Maped[currency+"1D"] = prices.Prices[0].Timestamp

	event1D := transformToEvent(prices, "1D", currency)

	err := k.PushPrices(event1D, "1D", currency)
	if err != nil {
		return nil, fmt.Errorf("could not send prices on 1D: %w", err)
	}
	return prices, nil
}

func sortPricesByTimeStamp(p *common.PriceResponseAPI) {
	sortField := func(field []common.PriceUnitAPI) {
		compare := func(i, j int) bool {
			return field[i].Timestamp > field[j].Timestamp
		}
		sort.Slice(field, compare)
	}

	sortField(p.Prices)
	sortField(p.MarketCaps)
	sortField(p.TotalVolumes)
}

func removeOlderValues(t int64, p *common.PriceResponseAPI) {
	removeOlder := func(field []common.PriceUnitAPI) []common.PriceUnitAPI {
		var result []common.PriceUnitAPI
		for _, entry := range field {
			if entry.Timestamp >= t {
				result = append(result, entry)
			}
		}
		return result
	}

	p.Prices = removeOlder(p.Prices)
	p.MarketCaps = removeOlder(p.MarketCaps)
	p.TotalVolumes = removeOlder(p.TotalVolumes)
}
