package zaif

import "net/url"

func (p *PrivateAPI) CancelOrder(currencyPair string, orderId string) (*TradeResponse, error) {
	params := url.Values{}
	if len(currencyPair) > 0 {
		params.Add("currency_pair", currencyPair)
	}
	params.Add("order_id", orderId)
	var r TradeResponse
	err := p.requestWithRetry("cancel_order", params, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
