package example

import (
	"fmt"

	"github.com/uphy/go-zaif/zaif"
)

func LastPrice() {
	publicApi := zaif.NewPublicAPI()
	price, _ := publicApi.LastPrice("xem_jpy")
	fmt.Printf("%#v\n", price)
}
