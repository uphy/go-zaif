package zaif

func (p *PublicAPI) CurrencyPairs() ([]string, error) {
	type Data struct {
		CurrencyPair string `json:"currency_pair"`
	}
	var data []Data
	if err := p.getWithRetry("currency_pairs/all", &data); err != nil {
		return nil, err
	}
	currencyPairs := []string{}
	for _, d := range data {
		currencyPairs = append(currencyPairs, d.CurrencyPair)
	}
	return currencyPairs, nil
}
