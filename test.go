package main

import (
	"fmt"
	// "sync"
	// "time"
)

func main() {

	if true {
		guessa()
		// fmt.Print(aa)
		return
	}

	// m := make(map[string]int)
	// var wg sync.WaitGroup
	// var lock sync.Mutex
	// wg.Add(2)

	// go func() {
	// 	for {
	// 		lock.Lock()
	// 		m["a"]++
	// 		lock.Unlock()
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		lock.Lock()
	// 		m["a"]++
	// 		fmt.Println(m["a"])
	// 		lock.Unlock()
	// 	}
	// }()

	// select {
	// case <-time.After(time.Second * 5):
	// 	fmt.Println("timeout, stopping")
	// }
}

func guessa() int {
	a := 1
	if true {
		defer func() {
			a = 22
			fmt.Println(a)
		}()

		fmt.Println(a)
		return a
	}
	return 0
}
