package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"

	pb "github.com/micnncim/ft/proto"
)

type client struct {
	cli  pb.FileTransfererClient
	conn *grpc.ClientConn
}

func newClient(host string) (*client, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := &client{
		cli:  pb.NewFileTransfererClient(conn),
		conn: conn,
	}
	return c, nil
}

func (c *client) run(ctx context.Context, cmd, from, to string) error {
	log.Printf("client: starting")

	errCh := make(chan error, 1)
	switch cmd {
	case "download":
		go func() {
			errCh <- c.download(ctx, from, to)
		}()
	case "upload":
		go func() {
			errCh <- c.upload(ctx, from, to)
		}()
	default:
		return errors.New("no such command")
	}

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return errors.New("client: canceled")
	}
}

func (c *client) download(ctx context.Context, from, to string) error {
	log.Printf("client: downloading %s...", from)

	f, err := os.Create(to)
	if err != nil {
		return err
	}

	for {
		stream, err := c.cli.Download(ctx, &pb.DownloadRequest{
			From: from,
		})
		if err != nil {
			return err
		}
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Printf("client: downloaded %s", to)
				return nil
			}
			if err != nil {
				return err
			}
			if _, err := f.Write(resp.Content); err != nil {
				return err
			}
		}
	}
}

func (c *client) upload(ctx context.Context, from, to string) error {
	log.Printf("client: uploading %s...", from)

	f, err := os.Open(from)
	if err != nil {
		return err
	}

	stream, err := c.cli.Upload(ctx)
	if err != nil {
		return err
	}

	b := make([]byte, 1024*4)
	for {
		n, err := f.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		req := &pb.UploadRequest{
			To:      to,
			Content: b[:n],
		}
		if err := stream.Send(req); err != nil {
			return err
		}
	}

	if _, err = stream.CloseAndRecv(); err != nil {
		return err
	}
	log.Printf("client: uploaded %s", to)
	return nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
