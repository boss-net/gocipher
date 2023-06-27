package main

import (
	"fmt"

	"github.com/boss-net/gocipher"
)

func main() {
	test := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "h"}
	bl := gocipher.New(10, 5)

	for i := 0; i < len(test); i++ {
		idx := bl.Shuffle(int64(i))
		fmt.Println(test[idx])
	}
}
