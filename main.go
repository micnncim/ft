package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"
)

var (
	serverMode bool
	port       string
	address    string

	timeout = 10 * time.Second
)

func init() {
	flag.BoolVar(&serverMode, "s", false, "run as a server")
	flag.StringVar(&port, "port", ":8080", "port for server")
	flag.StringVar(&address, "address", "localhost:8080", "address of server for client")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	if serverMode {
		f := &fileTransferer{}
		if err := f.Run(ctx); err != nil {
			log.Fatal(err)
		}
		return
	}

	if len(os.Args) != 4 {
		log.Fatal("Client command format is `[download|upload] <FILENAME_FROM> <FILENAME_TO>`")
	}
	cmd, from, to := os.Args[1], os.Args[2], os.Args[3]

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	c, err := newClient()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	if err := c.Run(ctx, cmd, from, to); err != nil {
		log.Fatal(err)
	}
}
