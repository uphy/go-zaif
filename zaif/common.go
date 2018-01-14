package zaif

import "errors"

type (
	Action string
	Price  float64
	Amount float64
	Order  struct {
		Price  Price
		Amount Amount
	}
	Depth struct {
		Asks []Order
		Bids []Order
	}
)

var (
	Err502          = errors.New("HTTP 502")
	ErrTooManyRetry = errors.New("too many retry")
)

const (
	ActionBid Action = "bid"
	ActionAsk Action = "ask"
	baseURL          = "https://api.zaif.jp/api/1/"
)

func newOrder(price Price, amount Amount) Order {
	return Order{price, amount}
}

func newDepth(asks []Order, bids []Order) *Depth {
	return &Depth{asks, bids}
}

func convertDepthArray(a [][]float64) []Order {
	var orders []Order
	for i := 0; i < len(a); i++ {
		orders = append(orders, newOrder(Price(a[i][0]), Amount(a[i][1])))
	}
	return orders
}
