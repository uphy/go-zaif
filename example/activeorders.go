package example

import "fmt"

func ActiveOrders() {
	privateAPI := newPrivateAPI()
	resp, err := privateAPI.ActiveOrders("")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
