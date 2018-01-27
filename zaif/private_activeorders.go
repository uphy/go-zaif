package zaif

import (
	"net/url"
	"time"
)

type (
	ActiveOrder struct {
		CurrencyPair string  `json:"currency_pair"`
		Action       string  `json:"action"`
		Amount       float64 `json:"amount"`
		Price        float64 `json:"price"`
		Timestamp    int64   `json:"timestamp,string"`
		Comment      string  `json:"comment"`
	}
	ActiveOrdersResponse map[string]ActiveOrder
)

func (a *ActiveOrder) Time() time.Time {
	return time.Unix(a.Timestamp, 0)
}

func (c *PrivateAPI) ActiveOrders(currencyPair string) (ActiveOrdersResponse, error) {
	params := url.Values{}
	if len(currencyPair) > 0 {
		params.Add("currency_pair", currencyPair)
	}
	var r ActiveOrdersResponse
	err := c.requestWithRetry("active_orders", params, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
