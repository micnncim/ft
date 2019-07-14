package main

import (
	"io"
	"os"

	pb "github.com/micnncim/ft/service"
)

type fileTransferer struct{}

func (f *fileTransferer) Download(req *pb.DownloadRequest, stream pb.FileTransferer_DownloadServer) error {
	file, err := os.Open(req.From)
	if err != nil {
		return err
	}

	buf := make([]byte, 1024*4)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
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
	file, err := os.Create(req.To)
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

	return stream.SendAndClose(&pb.UploadResponse{})
}
