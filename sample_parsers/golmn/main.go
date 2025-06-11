package main

import (
	"fmt"
	"lmn/lmn"
	"os"
	"reflect"
)

func main() {
	data, _ := os.ReadFile("test.lmn")
	s := string(data)
	val, err := lmn.Parse(s)

	if err != nil {
		println(err.Error())
	} else {
		fmt.Printf("%v\n", val)
		fmt.Printf("%T\n", val)
	}

	fmt.Println(reflect.TypeOf(any(nil)))
}
