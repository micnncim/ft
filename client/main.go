package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "github.com/micnncim/ft/proto"
)

const address = "localhost:8080"

func main() {
	if len(os.Args) != 4 {
		log.Fatal("Command format is `[download|upload] <FILENAME_FROM> <FILENAME_TO>`")
	}
	cmd, from, to := os.Args[1], os.Args[2], os.Args[3]

	// TODO: Use TLS.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := newClient(pb.NewFileTransfererClient(conn))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch cmd {
	case "download":
		log.Printf("Downloading %s...\n", from)
		if err := c.download(ctx, from, to); err != nil {
			log.Fatal(err)
		}
		log.Printf("Downloaded %s!\n", to)
	case "upload":
		log.Printf("Uploading %s...\n", from)
		if err := c.upload(ctx, from, to); err != nil {
			log.Fatal(err)
		}
		log.Printf("Uploaded %s!\n", to)
	default:
		log.Fatal("no such command")
	}
}
