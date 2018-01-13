package example

import "fmt"

func CancelOrder() {
	privateAPI := newPrivateAPI()
	resp, err := privateAPI.CancelOrder("", "300000000")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
