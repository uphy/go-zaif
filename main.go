package main

import "github.com/uphy/go-zaif/example"

func main() {
	go example.Stream()
	example.Chat()
}
