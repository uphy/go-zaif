package zaif

import (
	"sort"
	"time"
)

type zaifTrade struct {
	Date         int64
	Price        Price
	Amount       float64
	Tid          int64
	CurrencyPair string `json:"currency_pair"`
	TradeType    string `json:"trade_type"`
}

func (z *zaifTrade) convert() *Trade {
	return &Trade{
		Amount: Amount(z.Amount),
		Price:  z.Price,
		Time:   time.Unix(z.Date, 0),
		Action: Action(z.TradeType),
	}
}

func (p *PublicAPI) Trades(currencyPair string) ([]Trade, error) {
	var v []zaifTrade
	if err := p.getWithRetry("trades/"+currencyPair, &v); err != nil {
		return nil, err
	}
	sort.Slice(v, func(i, j int) bool {
		return v[i].Tid < v[j].Tid
	})
	var trades = make([]Trade, len(v))
	for i, t := range v {
		trades[i] = *t.convert()
	}
	return trades, nil
}
