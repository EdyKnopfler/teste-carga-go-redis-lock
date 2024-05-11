package main

import (
	"log"

	"golang.org/x/exp/constraints"
)

func sum[T constraints.Integer](a, b T) T {
	return a + b
}

func main() {
	var (
		a int   = 3
		b int   = 32
		c int32 = 3
		d int32 = 32
	)
	log.Println(sum(a, b))
	log.Println(sum(c, d))
}
