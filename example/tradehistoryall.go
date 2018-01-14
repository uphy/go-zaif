package example

import (
	"fmt"

	"github.com/uphy/go-zaif/zaif"
)

func TradeHistoryAll() {
	privateAPI := newPrivateAPI()
	c := make(chan zaif.TradeHistoryEntry, 10)
	e := make(chan error)
	privateAPI.TradeHistoryAll("xem_jpy", c, e)
l:
	for {
		select {
		case entry, ok := <-c:
			if ok {
				fmt.Println(entry)
			}
		case err, ok := <-e:
			if ok {
				fmt.Println(err)
			}
			break l
		}
	}
}
