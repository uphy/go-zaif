package example

import (
	"fmt"

	"github.com/uphy/go-zaif/zaif"
)

func Stream() {
	publicApi := zaif.NewPublicAPI()
	depth, trade, err := publicApi.Stream("xem_jpy")
l:
	for {
		select {
		case t := <-trade:
			fmt.Println(t)
		case d := <-depth:
			fmt.Println(d)
		case e := <-err:
			fmt.Println(e)
			break l
		}
	}
}
