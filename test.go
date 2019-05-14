package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {

	wg.Add(2)

	go doThis(10)
	go doThat(20)

	wg.Wait()
}

func doThis(i int) {
	defer wg.Done()
	fmt.Println("c : ", i)
}

func doThat(i int) {
	defer wg.Done()
	fmt.Println("c : ", i)
}
