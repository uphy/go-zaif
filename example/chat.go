package example

import (
	"fmt"

	"github.com/uphy/go-zaif/zaif"
)

func Chat() {
	publicAPI := zaif.NewPublicAPI()
	msg, err := publicAPI.Chat("/")
	for {
		select {
		case m := <-msg:
			fmt.Printf("%s %s\n", m.Time(), m.Message)
		case e := <-err:
			panic(e)
		}
	}
}
