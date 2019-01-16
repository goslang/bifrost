package main

import (
	"github.com/goslang/bifrost"
)

func main() {
	err := bifrost.Start()
	println(err.Error())
}
