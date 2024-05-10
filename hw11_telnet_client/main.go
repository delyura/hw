package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var ErrorArgumentsMissing = errors.New("Empty host or port") //nolint:all

func argumentParse() (address string, timeout time.Duration, err error) {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout to connect")
	flag.Parse()
	host := flag.Arg(0)
	port := flag.Arg(1)

	if host == "" || port == "" {
		return "", timeout, ErrorArgumentsMissing
	}
	return net.JoinHostPort(host, port), timeout, nil
}

func main() {
	address, timeout, err := argumentParse()
	if err != nil {
		log.Println(err)
	}

	ch := make(chan os.Signal, 1)
	telnetclient := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	ctx, cancel := context.WithCancel(context.Background())

	if err = telnetclient.Connect(); err != nil {
		log.Println("error Connect")
	}

	defer func(telnetclient TelnetClient) {
		if err := telnetclient.Close(); err != nil {
			log.Println("error Close")
		}
	}(telnetclient)

	go func() {
		defer cancel()
		if err := telnetclient.Receive(); err != nil {
			log.Println("error Receive")
			return
		}
	}()

	go func() {
		defer cancel()
		if err := telnetclient.Send(); err != nil {
			log.Println("error Send")
			return
		}
	}()

	signal.Notify(ch, os.Interrupt)
	select {
	case <-ch:
		cancel()
	case <-ctx.Done():
		log.Println("Connection was closed by peer")
		close(ch)
	}
}
