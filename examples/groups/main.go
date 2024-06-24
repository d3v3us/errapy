package main

import (
	"fmt"

	. "github.com/d3v3us/errapy/pkg"
)

var (
	_                    = Policy(WithClassesRequired(true), WithCodesRequired(true))
	Business             = Class("Business", "BE", nil)
	CustomerNotFound     = Business.New("customer not found")
	CurrencyDoesNotExist = Business.New("currency does not exist")
)

func main() {
	var err = Grouped(CustomerNotFound, CurrencyDoesNotExist)
	fmt.Println(err.Error())

	var group Group
	group.Add(CustomerNotFound)
	group.Add(CurrencyDoesNotExist)
	fmt.Println(group.Error())

	fmt.Printf("%+v\n", err)
}
