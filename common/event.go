package common

func CurrenciesToCurrencyEvents(currencies []Currency) ([]CurrencyEvent, error) {
	var currenciesEvent []CurrencyEvent

	for _, v := range currencies {
		cours := splitText(v.Cours, 1)
		var1h := splitText(v.Variation1h, 0)
		var1d := splitText(v.Variation1d, 0)
		var1w := splitText(v.Variation1w, 0)
		volume := splitText(v.Volume, 1)
		mkCap := splitText(v.MarketCapital, 1)
		volume24 := splitText(v.Volume24d, 1)
		fdv := splitText(v.FDV, 1)
		currenciesEvent = append(currenciesEvent,
			CurrencyEvent{
				ID:            stringToInt(v.ID),
				Name:          v.Name,
				Cours:         stringToFloat(cours),
				Variation1h:   stringToFloat(var1h),
				Variation1d:   stringToFloat(var1d),
				Variation1w:   stringToFloat(var1w),
				MarketCapital: stringToFloat(mkCap),
				Volume:        stringToFloat(volume),
				Chaine:        v.Chaine,
				Volume24d:     stringToFloat(volume24),
				FDV:           stringToFloat(fdv),
				LastAdded:     v.LastAdded,
			},
		)
	}

	return currenciesEvent, nil
}
