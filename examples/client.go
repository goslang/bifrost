package main

import (
	"flag"
	"fmt"
	"os"

	bifrost "github.com/goslang/bifrost/client"
)

func main() {
	channelName := flag.String("channel", "test", "Specify the channel to publish to.")
	message := flag.String("message", "", "The message to publish.")

	flag.Parse()

	cl, err := bifrost.New("127.0.0.1", 2727)
	exitOnErr(err)

	channel := bifrost.Channel{
		Name: *channelName,
	}

	err = cl.Do(channel.Publish(*message)).Error()
	exitOnErr(err)
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err.Error())
		os.Exit(1)
	}
}
