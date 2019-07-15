package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "github.com/micnncim/ft/proto"
)

type fileTransferer struct {
	host string
}

func newFileTransferer(host string) *fileTransferer {
	return &fileTransferer{
		host: host,
	}
}

func (f *fileTransferer) run(ctx context.Context) error {
	log.Printf("server: starting on %s", f.host)

	s := grpc.NewServer()
	pb.RegisterFileTransfererServer(s, f)

	lis, err := net.Listen("tcp", f.host)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Serve(lis)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return errors.New("client: timeout")
	}
}

func (f *fileTransferer) Download(req *pb.DownloadRequest, stream pb.FileTransferer_DownloadServer) error {
	from := req.From
	log.Printf("server: sending %s...", from)

	file, err := os.Open(from)
	if err != nil {
		return err
	}

	buf := make([]byte, 1024*4)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			log.Printf("server: sent %s", from)
			return nil
		}
		if err != nil {
			return err
		}
		err = stream.Send(&pb.DownloadResponse{
			Content: buf[:n],
		})
		if err != nil {
			return err
		}
	}
}

func (f *fileTransferer) Upload(stream pb.FileTransferer_UploadServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	to := req.To
	log.Printf("server: receiving %s...", to)
	file, err := os.Create(to)
	if err != nil {
		return err
	}

	for {
		if _, err := file.Write(req.Content); err != nil {
			return err
		}
		req, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	log.Printf("server: received %s", to)
	return stream.SendAndClose(&pb.UploadResponse{})
}
