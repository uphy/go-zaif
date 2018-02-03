package example

import (
	"fmt"
)

func Info() {
	privateAPI := newPrivateAPI()
	info, err := privateAPI.Info()
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
}
