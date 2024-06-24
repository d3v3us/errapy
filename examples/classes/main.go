package main

import (
	"fmt"

	. "github.com/d3v3us/errapy/pkg"
)

var (
	Authentication  = Class("Authentication", "BE", nil)
	ErrUnauthorized = Authentication.New("unauthorized user")
)

func main() {
	fmt.Println(ErrUnauthorized)
}
