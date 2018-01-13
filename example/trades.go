package example

import (
	"fmt"

	"github.com/uphy/go-zaif/zaif"
)

func Trades() {
	publicApi := zaif.NewPublicAPI()
	trades, _ := publicApi.Trades("xem_jpy")
	fmt.Printf("%#v\n", trades)
}
