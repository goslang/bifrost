package main

import (
	"flag"
	"fmt"
	"os"

	bifrost "github.com/goslang/bifrost/client"
)

func main() {
	channelName := flag.String("channel", "test", "Specify the channel to pop a message from")

	flag.Parse()

	cl, err := bifrost.New("127.0.0.1", 2727)
	exitOnErr(err)

	channel := bifrost.Channel{
		Name: *channelName,
	}

	var message string
	err = cl.Do(channel.Pop)(&message)
	exitOnErr(err)

	fmt.Println("Message:", message)
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err.Error())
		os.Exit(1)
	}
}
