package usecase

import (
	"fmt"
	"sort"

	"github.com/tobg8/crypto-viz/common"
	"github.com/tobg8/crypto-viz/internal/repository"
)

var last1DTimeStamp int64
var last90DTimeStamp int64
var lastALLimeStamp int64

func HandlePrices(k *repository.KafkaClient) error {
	// Je veux récupérer la currency et la range

	// Je veux vérifier que la cuurrency et la range existe

	// Je veux appeler ma func repo qui récup les prix en passant range et currency
	currency := "bitcoin"

	// Je dois faire 3 appels, 1 appel 1D, un appel 89 days, et un appel +90j
	_, err := handle1DPrice(k, currency)
	if err != nil {
		return err
	}
	_, err = handle90Price(k, currency)
	if err != nil {
		return err
	}
	_, err = handleAllPrice(k, currency)
	if err != nil {
		return err
	}

	return nil
}

func transformToEvent(prices *common.PriceResponseAPI, r string, currency string) []common.PriceEvent {
	var response []common.PriceEvent
	for _, v := range prices.Prices {
		if len(prices.Prices) > 0 {
			temp := common.PriceEvent{
				PriceUnitAPI: v,
				Type:         "prices",
				Range:        r,
				Currency:     currency,
			}
			response = append(response, temp)
		}
	}

	for _, v := range prices.TotalVolumes {
		if len(prices.Prices) > 0 {
			temp := common.PriceEvent{
				PriceUnitAPI: v,
				Type:         "total_volumes",
				Range:        r,
				Currency:     currency,
			}
			response = append(response, temp)
		}
	}

	for _, v := range prices.MarketCaps {
		if len(prices.Prices) > 0 {
			temp := common.PriceEvent{
				PriceUnitAPI: v,
				Type:         "market_caps",
				Range:        r,
				Currency:     currency,
			}
			response = append(response, temp)
		}
	}

	return response
}

func handleAllPrice(k *repository.KafkaClient, currency string) (*common.PriceResponseAPI, error) {
	prices := k.FetchPrices(currency, "300")
	if prices == nil {
		return nil, fmt.Errorf("error when fetching prices ALL")
	}
	// Je veux trier mes arrays par timestamp,
	sortPricesByTimeStamp(prices)

	if lastALLimeStamp == prices.Prices[0].Timestamp {
		return nil, fmt.Errorf("no new prices to send on ALL")
	}
	// J'enlève ceux qui sont plus grands
	removeOlderValues(lastALLimeStamp, prices)

	// Je set le lastTimeStamp avec le plus récent
	lastALLimeStamp = prices.Prices[0].Timestamp

	eventALL := transformToEvent(prices, "ALL", currency)
	err := k.PushPrices(eventALL, "ALL")
	if err != nil {
		return nil, fmt.Errorf("could not send prices on ALL: %w", err)
	}
	return prices, nil
}

func handle1DPrice(k *repository.KafkaClient, currency string) (*common.PriceResponseAPI, error) {
	prices := k.FetchPrices(currency, "1")
	if prices == nil {
		return nil, fmt.Errorf("error when fetching prices on 1D")
	}
	// Je veux trier mes arrays par timestamp,
	sortPricesByTimeStamp(prices)

	if last1DTimeStamp == prices.Prices[0].Timestamp {
		return nil, fmt.Errorf("no new prices to send on 1D")
	}
	// J'enlève ceux qui sont plus grands
	removeOlderValues(last1DTimeStamp, prices)

	// Je set le lastTimeStamp avec le plus récent
	last1DTimeStamp = prices.Prices[0].Timestamp

	event1D := transformToEvent(prices, "1D", currency)

	err := k.PushPrices(event1D, "1D")
	if err != nil {
		return nil, fmt.Errorf("could not send prices on 1D: %w", err)
	}
	return prices, nil
}

func handle90Price(k *repository.KafkaClient, currency string) (*common.PriceResponseAPI, error) {
	prices := k.FetchPrices(currency, "89")
	if prices == nil {
		return nil, fmt.Errorf("error when fetching prices on 90D")
	}
	// Je veux trier mes arrays par timestamp,
	sortPricesByTimeStamp(prices)

	if last90DTimeStamp == prices.Prices[0].Timestamp {
		return nil, fmt.Errorf("no new prices to send on 90D")
	}
	// J'enlève ceux qui sont plus grands
	removeOlderValues(last90DTimeStamp, prices)

	// Je set le lastTimeStamp avec le plus récent
	last90DTimeStamp = prices.Prices[0].Timestamp
	event90D := transformToEvent(prices, "90D", currency)
	err := k.PushPrices(event90D, "90D")
	if err != nil {
		return nil, fmt.Errorf("could not send prices on 90D: %w", err)
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
