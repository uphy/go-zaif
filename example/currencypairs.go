package example

import (
	"fmt"

	"github.com/uphy/go-zaif/zaif"
)

func CurrencyPairs() {
	publicApi := zaif.NewPublicAPI()
	c, _ := publicApi.CurrencyPairs()
	fmt.Printf("%#v\n", c)
}
