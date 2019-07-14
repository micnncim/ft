package main

import (
	"context"
	"io"
	"os"

	pb "github.com/micnncim/ft/proto"
)

type client struct {
	cli pb.FileTransfererClient
}

func newClient(fileTransfererClient pb.FileTransfererClient) *client {
	return &client{
		cli: fileTransfererClient,
	}
}

func (c *client) download(ctx context.Context, from, to string) error {
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
	return nil
}
