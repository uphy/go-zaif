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
	Board struct {
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

func newBoard(asks []Order, bids []Order) *Board {
	return &Board{asks, bids}
}

func convertBoardArray(a [][]float64) []Order {
	var orders []Order
	for i := 0; i < len(a); i++ {
		orders = append(orders, newOrder(Price(a[i][0]), Amount(a[i][1])))
	}
	return orders
}
