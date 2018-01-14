package example

import (
	"fmt"
	"time"
)

func TradeHistory() {
	privateAPI := newPrivateAPI()
	to := time.Now()
	from := to.Add(-time.Hour * 24)
	tradeHistory, err := privateAPI.TradeHistoryByPeriod("xem_jpy", &from, &to)
	if err != nil {
		panic(err)
	}
	fmt.Println(tradeHistory)
}
