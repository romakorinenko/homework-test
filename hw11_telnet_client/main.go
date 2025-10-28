package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeoutStr string
	flag.StringVar(&timeoutStr, "timeout", "10s", "Duration of connections")
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatal("required arguments \"host\" and \"port\" not define")
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Fatalf("timeout is invalid: %s", err)
	}

	telnetClient := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err = telnetClient.Connect(); err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
	defer func() {
		if err = telnetClient.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := telnetClient.Send(); err != nil {
			log.Println(err)
			return
		}
	}()

	go func() {
		defer close(done)
		if err := telnetClient.Receive(); err != nil {
			log.Println(err)
			return
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-done:
		fmt.Fprintln(os.Stderr, "...Connection was closed by peer or EOF")
	case <-signalChan:
		fmt.Fprintln(os.Stderr, "...Interrupted")
		close(signalChan)
	}
}
