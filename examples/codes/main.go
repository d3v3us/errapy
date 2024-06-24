package main

import (
	"fmt"

	. "github.com/d3v3us/errapy/pkg"
)

var (
	_                    = Policy(WithClassesRequired(true), WithCodesRequired(true))
	Business             = Class("Business", "BE", &[]int{100, 200})
	CustomerNotFound     = Business.New("customer not found")
	CurrencyDoesNotExist = Business.New("currency does not exist")
)

func main() {
	fmt.Println(CustomerNotFound.Error())
	fmt.Println(CurrencyDoesNotExist.Error())
}
