package example

import (
	"fmt"

	"github.com/uphy/go-zaif/zaif"
)

func Board() {
	publicApi := zaif.NewPublicAPI()
	board, _ := publicApi.GetBoard("xem_jpy")
	fmt.Printf("%#v\n", board.Asks)
	fmt.Printf("%#v\n", board.Bids)
}
