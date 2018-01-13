package zaif

import (
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
		Board     *Board
		LastPrice LastPrice
		Trades    []Trade
		Time      time.Time
	}
)

func (p *PublicAPI) Stream(currencyPair string, b chan<- Board, t chan<- Trade, err chan<- error) {
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
		var lastTid int64 = -1
		for {
			var v Data
			if e := websocket.JSON.Receive(ws, &v); e != nil {
				err <- e
			}
			b <- *newBoard(convertBoardArray(v.Asks), convertBoardArray(v.Bids))
			for _, trade := range v.Trades {
				if trade.Tid > lastTid {
					t <- Trade{
						Action: Action(trade.TradeType),
						Price:  Price(trade.Price),
						Amount: Amount(trade.Amount),
						Time:   time.Unix(trade.Date, 0),
					}
					lastTid = trade.Tid
				}
			}
		}
	}()
}
