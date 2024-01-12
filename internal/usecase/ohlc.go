package usecase

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/tobg8/crypto-viz/common"
	"github.com/tobg8/crypto-viz/internal/repository"
)

var OhlcMap map[string]int64

func InitMapOhlc() {
	// Initialize the map before using it
	OhlcMap = make(map[string]int64)
}

func HandleOhlc(k *repository.KafkaClient) error {
	log.Print("new ohlc")
	// Je veux appeler ma func repo qui récup les prix en passant range et currency
	currencies := repository.FetchListing()
	// Je les sauvegarde
	for _, v := range *currencies {
		if _, ok := OhlcMap[v.ID+"1D"]; !ok {
			OhlcMap[v.ID+"1D"] = 0
		}
	}

	for i, v := range *currencies {
		// Pour l'instant on gère 10 currencies en One day
		// Parce que on prends de temps que les 5 mins de la routine
		if i == 11 {
			return nil
		}
		_, err := handleOhlc(k, v.ID)
		if err != nil {
			log.Print(err)
			// here we must wait becasue coingecko rate limit returns 401
			time.Sleep(time.Second * 13)
			continue
		}
	}

	return nil
}

func handleOhlc(k *repository.KafkaClient, currency string) (*common.OhlcResponseAPI, error) {
	graphs := k.FetchOhlc(currency, "1")
	if graphs == nil {
		return nil, fmt.Errorf("error when fetching ohlc on 1D")
	}
	// Je veux trier mes arrays par timestamp,
	sortOhlcByTimeStamp(graphs)

	// Check if the stored timestamp for the currency is equal to or older than the latest timestamp
	if storedTimestamp, ok := OhlcMap[currency+"1D"]; ok && storedTimestamp >= graphs.Ohlc[0].Timestamp {
		return nil, fmt.Errorf("no new ohlc to send on 1D")
	}

	// J'enlève ceux qui sont plus grands
	removeOlderOhlcValues(OhlcMap[currency+"1D"], graphs)

	// Je set le lastTimeStamp avec le plus récent
	OhlcMap[currency+"1D"] = graphs.Ohlc[0].Timestamp

	event1D := transformToOhlc(graphs, "1D", currency)

	err := k.PushOhlc(event1D, "1D", currency)
	if err != nil {
		return nil, fmt.Errorf("could not send prices on 1D: %w", err)
	}
	return graphs, nil
}

func transformToOhlc(graphs *common.OhlcResponseAPI, r string, currency string) common.OhlcResponseTest {
	var response common.OhlcResponseTest
	for _, v := range graphs.Ohlc {
		if len(graphs.Ohlc) > 0 {
			temp := common.OhlcEvent{
				OhlcUnit: v,
				Currency: currency,
				Range:    r,
			}
			response.Ohlc = append(response.Ohlc, temp)
		}
	}
	response.Currency = currency
	response.Range = r
	return response
}

func sortOhlcByTimeStamp(p *common.OhlcResponseAPI) {
	sortField := func(field []common.OhlcUnit) {
		compare := func(i, j int) bool {
			return field[i].Timestamp > field[j].Timestamp
		}
		sort.Slice(field, compare)
	}

	sortField(p.Ohlc)
}

func removeOlderOhlcValues(t int64, p *common.OhlcResponseAPI) {
	removeOlder := func(field []common.OhlcUnit) []common.OhlcUnit {
		var result []common.OhlcUnit
		for _, entry := range field {
			if entry.Timestamp >= t {
				result = append(result, entry)
			}
		}
		return result
	}

	p.Ohlc = removeOlder(p.Ohlc)

}
