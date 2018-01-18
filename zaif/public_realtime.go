package zaif

import (
	"sort"
	"time"

	"golang.org/x/net/websocket"
)

type (
	Trade struct {
		Action Action
		Price  Price
		Amount Amount
		Time   time.Time
	}

	LastPrice struct {
		Action Action
		Price  Price
	}

	StreamData struct {
		Depth     *Depth
		LastPrice LastPrice
		Trades    []Trade
		Time      time.Time
	}
)

func (p *PublicAPI) Stream(currencyPair string) (<-chan Depth, <-chan Trade, <-chan error) {
	depth := make(chan Depth, 100)
	trade := make(chan Trade, 100)
	err := make(chan error, 10)
	type RawTrade struct {
		CurrencyPair string  `json:"currency_pair"`
		TradeType    string  `json:"trade_type"`
		Price        float64 `json:"price"`
		Amount       float64 `json:"amount"`
		Date         int64   `json:"date"`
		Tid          int64   `json:"tid"`
	}

	type Data struct {
		Asks         [][]float64 `json:"asks"`
		Bids         [][]float64 `json:"bids"`
		Trades       []RawTrade  `json:"trades"`
		Timestamp    string      `json:"timestamp"`
		LastPrice    LastPrice   `json:"last_price"`
		CurrencyPair string      `json:"currency_pair"`
	}

	ws, e := websocket.Dial("wss://ws.zaif.jp:8888/stream?currency_pair="+currencyPair, "", "https://ws.zaif.jp:8888")
	if e != nil {
		err <- e
	}
	go func() {
		defer close(depth)
		defer close(trade)
		defer close(err)
		var lastTid int64 = -1
		for {
			var v Data
			if e := websocket.JSON.Receive(ws, &v); e != nil {
				err <- e
			}
			depth <- *newDepth(convertDepthArray(v.Asks), convertDepthArray(v.Bids))
			sort.Slice(v.Trades, func(i, j int) bool {
				return v.Trades[i].Tid < v.Trades[j].Tid
			})
			for _, t := range v.Trades {
				if t.Tid > lastTid {
					trade <- Trade{
						Action: Action(t.TradeType),
						Price:  Price(t.Price),
						Amount: Amount(t.Amount),
						Time:   time.Unix(t.Date, 0),
					}
					lastTid = t.Tid
				}
			}
		}
	}()
	return depth, trade, err
}
