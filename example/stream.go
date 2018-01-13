package example

import (
	"fmt"

	"github.com/uphy/go-zaif/zaif"
)

func Stream() {
	publicApi := zaif.NewPublicAPI()

	boardChannel := make(chan zaif.Board, 10)
	tradeChannel := make(chan zaif.Trade, 10)
	errorChannel := make(chan error, 10)
	publicApi.Stream("xem_jpy", boardChannel, tradeChannel, errorChannel)
l:
	for {
		select {
		case trade := <-tradeChannel:
			fmt.Println(trade)
		case board := <-boardChannel:
			fmt.Println(board)
		case err := <-errorChannel:
			fmt.Println(err)
			break l
		}
	}
}
