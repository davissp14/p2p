package server

import (
	"context"
	"io"
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
	log.Println("Initiating file transfer...")
	buf := make([]byte, 1024)
	writing := true
	success := false
	for writing {
		n, err := file.Read(buf)
		if err == io.EOF {
			writing = false
			success = true
			break
		}
		if err != nil {
			log.Printf("failed to read file. error: %s", err.Error())
			break
		}
		chunk := pb.Chunk{Data: buf[:n]}
		stream.Send(&chunk)
	}

	if success {
		log.Println("File transfer completed!")
	}

	return nil
}
