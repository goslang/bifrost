package main

import (
	"flag"
	"fmt"
	"os"

	bifrost "github.com/goslang/bifrost/client"
)

func main() {
	name := flag.String("channel", "", "The name of the channel to remove.")
	flag.Parse()

	cl, err := bifrost.New("127.0.0.1", 2727)
	exitOnErr(err)

	err = cl.Do(bifrost.DestroyChannel(*name)).Error()
	exitOnErr(err)

	fmt.Println("Removed channel", *name)
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err.Error())
		os.Exit(1)
	}
}
