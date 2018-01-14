package zaif

import (
	"fmt"
	"net/url"
)

type TradeParameter struct {
	CurrencyPair string
	Action       Action
	Price        Price
	Amount       Amount
	Limit        Price
	Comment      string
}

type TradeResponse struct {
	Received float64
	Remains  float64
	OrderID  int64
	Funds    map[string]float64
}

func NewTradeParameter(currencyPair string, action Action, price Price, amount Amount) *TradeParameter {
	return &TradeParameter{
		CurrencyPair: currencyPair,
		Action:       action,
		Price:        price,
		Amount:       amount,
		Limit:        -1,
		Comment:      "",
	}
}

func (c *PrivateAPI) Trade(p *TradeParameter) (*TradeResponse, error) {
	params := url.Values{}
	params.Add("currency_pair", p.CurrencyPair)
	params.Add("action", string(p.Action))
	params.Add("price", fmt.Sprint(p.Price))
	params.Add("amount", fmt.Sprint(p.Amount))
	if p.Limit > 0 {
		params.Add("limit", fmt.Sprint(p.Limit))
	}
	if len(p.Comment) > 0 {
		params.Add("comment", p.Comment)
	}
	var data TradeResponse
	if err := c.requestWithRetry("trade", params, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
