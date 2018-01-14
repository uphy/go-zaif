package example

import (
	"fmt"

	"github.com/uphy/go-zaif/zaif"
)

func Stream() {
	publicApi := zaif.NewPublicAPI()

	depthChannel := make(chan zaif.Depth, 10)
	tradeChannel := make(chan zaif.Trade, 10)
	errorChannel := make(chan error, 10)
	publicApi.Stream("xem_jpy", depthChannel, tradeChannel, errorChannel)
l:
	for {
		select {
		case trade := <-tradeChannel:
			fmt.Println(trade)
		case depth := <-depthChannel:
			fmt.Println(depth)
		case err := <-errorChannel:
			fmt.Println(err)
			break l
		}
	}
}
