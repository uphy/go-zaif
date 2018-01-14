package example

import (
	"fmt"

	"github.com/uphy/go-zaif/zaif"
)

func Depth() {
	publicApi := zaif.NewPublicAPI()
	depth, _ := publicApi.GetDepth("xem_jpy")
	fmt.Printf("%#v\n", depth.Asks)
	fmt.Printf("%#v\n", depth.Bids)
}
