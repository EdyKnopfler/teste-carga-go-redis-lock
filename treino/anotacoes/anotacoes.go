package main

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Gato struct {
	Nome  string `transform:"pontinhos"`
	Raca  string
	Cor   string
	Idade uint
}

func StringComPontinhos(s string) string {
	var buffer bytes.Buffer

	for _, c := range s {
		buffer.WriteString(string(c))
		buffer.WriteString(".")
	}

	return buffer.String()
}

func Analisa(obj any) error {
	objType := reflect.TypeOf(obj).Elem()
	objValue := reflect.ValueOf(obj).Elem()

	if objValue.Kind() != reflect.Struct {
		return errors.New("Queremos uma struct!")
	}

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		value := objValue.Field(i)

		if !value.CanSet() {
			continue
		}

		transform := field.Tag.Get("transform")
		if transform == "" || transform == "-" {
			continue
		}

		if strings.Contains(transform, "pontinhos") {
			if value.Kind() == reflect.String {
				value.SetString(StringComPontinhos(value.String()))
			} else {
				return errors.New("Tipo não suportado")
			}
		}
	}

	return nil
}

func main() {
	kania := Gato{
		Nome:  "Scania",
		Raca:  "viralata",
		Cor:   "preto e branco",
		Idade: 11,
	}

	if err := Analisa(&kania); err != nil {
		panic(err)
	}

	fmt.Println(kania)
}
