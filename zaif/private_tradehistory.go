package zaif

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type (
	TradeHistory      []TradeHistoryEntry
	TradeHistoryEntry struct {
		ID           int64
		CurrencyPair string `json:"currency_pair"`
		Action       Action
		Amount       Amount
		Price        Price
		Fee          Amount
		YourAction   Action `json:"your_action"`
		Bonus        Amount
		Timestamp    int64 `json:",string"`
		Comment      string
	}
)

func (t *TradeHistoryEntry) ComputeProfit() Price {
	profit := 0.
	switch t.YourAction {
	case ActionBid:
		profit -= float64(t.Price) * float64(t.Amount)
	case ActionAsk:
		profit += float64(t.Price) * float64(t.Amount)
	}
	profit -= float64(t.Fee) * float64(t.Price)
	if t.Bonus > 0 {
		profit += float64(t.Bonus)
	}
	return Price(profit)
}

func (p *PrivateAPI) TradeHistoryByID(currencyPair string, id int64) (*TradeHistoryEntry, error) {
	tradeHistory, err := p.TradeHistoryByIDInterval(currencyPair, &id, &id)
	if err != nil {
		return nil, err
	}
	switch len(tradeHistory) {
	case 0:
		return nil, nil
	case 1:
		for _, entry := range tradeHistory {
			return &entry, nil
		}
		panic("unexpected response")
	default:
		return nil, errors.New("unexpected response")
	}
}

func (p *PrivateAPI) TradeHistoryByPeriod(currencyPair string, from *time.Time, to *time.Time) (TradeHistory, error) {
	params := url.Values{}
	if from != nil {
		params.Add("since", fmt.Sprint(from.Unix()))
	}
	if to != nil {
		params.Add("end", fmt.Sprint(to.Unix()))
	}
	return p.tradeHistory(currencyPair, params)
}

func (p *PrivateAPI) TradeHistoryByIDInterval(currencyPair string, fromID *int64, toID *int64) (TradeHistory, error) {
	params := url.Values{}
	if fromID != nil {
		params.Add("from_id", fmt.Sprint(*fromID))
	}
	if toID != nil {
		params.Add("end_id", fmt.Sprint(*toID))
	}
	return p.tradeHistory(currencyPair, params)
}

func (p *PrivateAPI) TradeHistoryAll(currencyPair string, c chan<- TradeHistoryEntry, e chan<- error) {
	go func() {
		defer close(c)
		defer close(e)

		from := 0
		count := 500
		for {
			tradeHistory, err := p.TradeHistoryByIndex(currencyPair, &from, &count)
			if err != nil {
				e <- err
				break
			}
			if len(tradeHistory) == 0 {
				break
			}
			for _, entry := range tradeHistory {
				c <- entry
			}
			from += count
		}
	}()
}

func (p *PrivateAPI) TradeHistoryByIndex(currencyPair string, from *int, count *int) (TradeHistory, error) {
	params := url.Values{}
	if from != nil {
		params.Add("from", fmt.Sprint(*from))
	}
	if count != nil {
		params.Add("count", fmt.Sprint(*count))
	}
	return p.tradeHistory(currencyPair, params)
}

func (p *PrivateAPI) tradeHistory(currencyPair string, params url.Values) (TradeHistory, error) {
	params.Add("currency_pair", currencyPair)

	var v map[string]TradeHistoryEntry
	if err := p.requestWithRetry("trade_history", params, &v); err != nil {
		return nil, err
	}
	var history TradeHistory
	for id, entry := range v {
		idint, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		entry.ID = idint
		history = append(history, entry)
	}
	sort.Slice(history, func(i, j int) bool {
		return history[i].ID < history[j].ID
	})
	return history, nil
}
