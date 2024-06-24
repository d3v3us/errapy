package main

import (
	"fmt"

	. "github.com/d3v3us/errapy/pkg"
)

var (
	_              = Policy(WithClassesRequired(true), WithCodesRequired(true))
	Authentication = Class("Authentication", "SE", nil)
	Unauthorized   = Authentication.New("unauthorized user")
)

func main() {
	fmt.Println(Unauthorized)
}
