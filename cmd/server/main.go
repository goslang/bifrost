package main

import (
	"github.com/goslang/bifrost/server"
)

func main() {
	err := server.Start()
	println(err.Error())
}
