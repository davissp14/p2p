package server

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/davissp14/p2p/pkg/service"
)

type PeerServiceServer struct {
	tls      bool
	keyFile  string
	certFile string
	port     int
}

func NewServer(port int, tls bool, keyFile, certFile string) *PeerServiceServer {
	return &PeerServiceServer{
		port:     port,
		tls:      tls,
		certFile: certFile,
		keyFile:  keyFile,
	}
}

func (p *PeerServiceServer) Ping(ctx context.Context, in *empty.Empty) (*pb.PingMessage, error) {
	response := pb.PingMessage{Message: "PONG"}
	return &response, nil
}

func (p *PeerServiceServer) Download(req *pb.PeerDownloadRequest, stream pb.PeerService_DownloadServer) error {
	file, err := os.Open(req.FilePath)
	if err != nil {
		return err
	}
	log.Printf("Incoming request to download file `%s`\n", req.FilePath)
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		chunk := pb.Chunk{Data: buf[:n]}
		stream.Send(&chunk)
	}
	log.Println("File transfer completed!")
	return nil
}

func (p *PeerServiceServer) List(req *pb.ListRequest, stream pb.PeerService_ListServer) error {
	fs, err := ioutil.ReadDir(req.Directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range fs {
		file := pb.File{
			Name:  f.Name(),
			Size:  f.Size(),
			IsDir: f.IsDir(),
		}
		stream.Send(&file)
	}

	return nil
}
