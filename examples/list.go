package main

import (
	"fmt"
	"os"

	bifrost "github.com/goslang/bifrost/client"
)

func main() {
	cl, err := bifrost.New("127.0.0.1", 2727)
	exitOnErr(err)

	var channels []bifrost.Channel
	err = cl.Do(bifrost.ListChannels)(&channels)
	exitOnErr(err)

	fmt.Println("Name\t\tMax\tSize")
	for _, c := range channels {
		fmt.Printf("%v\t%v\t%v\n", c.Name, c.Max, c.Size)
	}
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err.Error())
		os.Exit(1)
	}
}
