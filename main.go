package main

import (
	"chatbox/internal"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	port       = flag.Int("port", 50051, "The server port")
	clientMode = flag.Bool("client", false, "Chose to run as server, providing no flat forces client mode.")
	serverAddr = flag.String("hostname", "localhost", "The server host name")
	userflag   = flag.String("username", "test-user", "Choose a unique name")
)

func main() {
	flag.Parse()
	host := fmt.Sprintf("%s:%d", string(*serverAddr), int(*port))

	if *clientMode {
		client := internal.Client(host, *userflag)
		err := client.Run(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}

	server, listener := internal.Server(port)
	server.Serve(listener)
}
