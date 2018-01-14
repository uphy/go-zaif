package zaif

func (p *PublicAPI) LastPrice(currencyPair string) (Price, error) {
	var v map[string]float64
	if err := p.getWithRetry("last_price/"+currencyPair, &v); err != nil {
		return 0, err
	}
	return Price(v["last_price"]), nil
}
