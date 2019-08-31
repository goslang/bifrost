package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"

	bifrost "github.com/goslang/bifrost/client"
)

var cli = struct {
	channel struct {
		name string
		size uint
	}

	message struct {
		buf string
	}
}{}

func main() {
	cl, err := bifrost.New("127.0.0.1", 2727)
	exitOnErr(err)

	args := os.Args

	if len(args) < 2 {
		exitOnErr(errors.New("Must specify a command!"))
	}
	cmd := args[1]
	os.Args = args[1:]

	switch cmd {
	case "create":
		flag.UintVar(&cli.channel.size, "size", 5, "The size of the channel.")
		flag.StringVar(&cli.channel.name, "channel", "", "The name of the channel.")
	case "push":
		flag.StringVar(&cli.message.buf, "message", "", "The message to push.")
		flag.StringVar(&cli.channel.name, "channel", "", "The name of the channel.")
	default:
		flag.StringVar(&cli.channel.name, "channel", "", "The name of the channel.")
	}
	flag.Parse()

	switch cmd {
	case "create":
		err = Create(cl)
	case "list":
		err = List(cl)
	case "push":
		err = Push(cl)
	case "pop":
		err = Pop(cl)
	case "remove":
		err = Remove(cl)
	}
	exitOnErr(err)
}

func Create(cl *bifrost.Client) error {
	err := cl.Do(
		bifrost.CreateChannel(cli.channel.name, cli.channel.size),
	).Error()

	if err != nil {
		return err
	}
	return nil
}

func List(cl *bifrost.Client) error {
	var channels []bifrost.Channel
	if err := cl.Do(bifrost.ListChannels)(&channels); err != nil {
		return err
	}

	fmt.Println("Name\t\tMax\tSize")
	for _, c := range channels {
		fmt.Printf("%v\t%v\t%v\n", c.Name, c.Max, c.Size)
	}
	return nil
}

func Pop(cl *bifrost.Client) error {
	channel := bifrost.Channel{
		Name: cli.channel.name,
	}

	var message string
	if err := cl.Do(channel.Pop)(&message); err != nil {
		return err
	}

	fmt.Println("Message:", message)
	return nil
}

func Push(cl *bifrost.Client) error {
	channel := bifrost.Channel{
		Name: cli.channel.name,
	}

	err := cl.Do(channel.Publish(cli.message.buf)).Error()
	if err != nil {
		return err
	}
	return nil
}

func Remove(cl *bifrost.Client) error {
	err := cl.Do(
		bifrost.DestroyChannel(cli.channel.name),
	).Error()

	if err != nil {
		return err
	}
	return nil
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err.Error())
		os.Exit(1)
	}
}
