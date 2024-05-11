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

func main() {
	x := Coiso{Nome: "Kânia"}
	fmt.Println(x.AgoraVai())
}
