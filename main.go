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
	host       string

	timeout = 10 * time.Second
)

func init() {
	flag.BoolVar(&serverMode, "s", false, "run as a server")
	flag.StringVar(&host, "host", "localhost:8080", "host of server")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	if serverMode {
		f := newFileTransferer(host)
		if err := f.run(ctx); err != nil {
			log.Fatal(err)
		}
		return
	}

	if len(os.Args) != 4 {
		log.Fatal("client command format is `[download|upload] <FILENAME_FROM> <FILENAME_TO>`")
	}
	cmd, from, to := os.Args[1], os.Args[2], os.Args[3]

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	c, err := newClient(host)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	if err := c.run(ctx, cmd, from, to); err != nil {
		log.Fatal(err)
	}
}
