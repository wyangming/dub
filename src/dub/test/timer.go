package main

import (
	"fmt"
	"time"
)

func main() {
	/*c := time.Tick(1 * time.Second)
	for now := range c {
		fmt.Printf("%v aaa\n", now)
	}*/
	timer := time.NewTicker(2 * time.Second)
	for _ = range timer.C {
		fmt.Println(fmt.Sprintf("eee "))
	}
}
