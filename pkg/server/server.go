package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/davissp14/p2p/pkg/service"
	pb "github.com/davissp14/p2p/pkg/service"

	"google.golang.org/grpc"
)

var (
	hostname = flag.String("addr", "localhost", "The server hostname")
	port     = flag.Int("port", 8080, "The server port")
)

type peerServiceServer struct {
	mu sync.Mutex
}

func (e *peerServiceServer) Ping(ctx context.Context, msg *service.PingMessage) (*service.PingMessage, error) {
	response := pb.PingMessage{Message: "PONG"}
	return &response, nil
}

func newServer() *peerServiceServer {
	return &peerServiceServer{}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *hostname, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Listening on port %d\n", *port)
	var opts []grpc.ServerOption
	// opts = []grpc.ServerOption{grpc.WithInsecure()}
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterPeerServiceServer(grpcServer, newServer())

	grpcServer.Serve(lis)
}
