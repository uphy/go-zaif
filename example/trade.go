package example

import (
	"github.com/uphy/go-zaif/zaif"
)

func Trade() {
	privateAPI := newPrivateAPI()
	if _, err := privateAPI.Trade(zaif.NewTradeParameter("xem_jpy", zaif.ActionBid, 114.9999, 3106)); err != nil {
		panic(err)
	}
}
