package example

import "fmt"

func Info2() {
	privateAPI := newPrivateAPI()
	info2, err := privateAPI.Info2()
	if err != nil {
		panic(err)
	}
	fmt.Println(info2)
}
