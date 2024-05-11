package main

import (
	"fmt"
	"os"
	"runtime"
)

func soma(dados []uint8, ch chan uint64) {
	total := uint64(0)

	for _, theByte := range dados {
		total += uint64(theByte)
	}

	ch <- total
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		panic("Informe o arquivo!")
	}

	dados, err := os.ReadFile(args[0])
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%T\n", dados) // []uint8

	quantas := runtime.NumCPU() - 1
	tamanhoPedaco := len(dados) / quantas
	distribuido := 0
	ch := make(chan uint64, quantas)

	for i := 0; i < quantas; i++ {
		go soma(dados[i*tamanhoPedaco:i*tamanhoPedaco+tamanhoPedaco], ch)
		distribuido += tamanhoPedaco
	}

	if distribuido < len(dados) {
		quantas++
		go soma(dados[distribuido:len(dados)], ch)
	}

	total := uint64(0)

	for i := 0; i < quantas; i++ {
		total += <-ch
	}

	fmt.Println(total)
}
