package main

import "fmt"

type Teste interface {
	AgoraVai() string
}

type Coiso struct {
	Nome string
}

func (c *Coiso) AgoraVai() string {
	return fmt.Sprintf("O nome é: %s", c.Nome)
}

func Imprimir(t Teste) {
	fmt.Println(t.AgoraVai())
}

func main() {
	x := Coiso{Nome: "Kânia"}
	Imprimir(&x)
}
