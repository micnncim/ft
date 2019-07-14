package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/micnncim/ft/service"
)

func main() {
	const port = ":8080"

	s := grpc.NewServer()
	pb.RegisterFileTransfererServer(s, &fileTransferer{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
