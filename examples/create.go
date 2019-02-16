package main

import (
	"flag"
	"fmt"
	"os"

	bifrost "github.com/goslang/bifrost/client"
)

func main() {
	name := flag.String("channel", "", "The name of the new channel.")
	size := flag.Uint("size", 5, "The size of the new channel.")
	flag.Parse()

	cl, err := bifrost.New("127.0.0.1", 2727)
	exitOnErr(err)

	err = cl.Do(bifrost.CreateChannel(*name, *size)).Error()
	exitOnErr(err)

	fmt.Println("Created new channel")
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err.Error())
		os.Exit(1)
	}
}
