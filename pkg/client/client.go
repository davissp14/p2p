package client

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/davissp14/p2p/pkg/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	targetAddr = flag.String("targetAddr", "localhost:8080", "endpoint of your service")
	certFile   = flag.String("cert-file", "", "identify HTTPS client using this SSL certificate file")
)

func Ping(client pb.PeerServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Ping(ctx, &empty.Empty{})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resp.GetMessage())
}

func Download(client pb.PeerServiceClient, peerAddr, filePath string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := pb.PeerDownloadRequest{
		Addr:     peerAddr,
		FilePath: filePath,
	}
	stream, err := client.Download(ctx, &req)
	if err != nil {
		log.Fatalf("failed to initiate stream. error: %s", err.Error())
	}

	name := filepath.Base(filePath)

	file, err := os.Create(fmt.Sprintf("downloads/%s", name))
	if err != nil {
		log.Printf("Failed to create cert file: error. %s", err.Error())
	}

	defer file.Close()
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err.Error())
			break
		}
		_, err = file.Write(in.Data)
		if err != nil {
			log.Printf("Failed to stream file. error: %s", err.Error())
			break
		}
	}
	log.Printf("File `%s` download was a success!\n", name)
}

func NewClient(cert, addr, filePath string) {

	// Default to insecure connection
	transportStrategy := grpc.WithInsecure()
	// Establish Secure connection using provided public key
	if cert != "" {
		creds, err := credentials.NewClientTLSFromFile(cert, "")
		if err != nil {
			fmt.Println("missing or invalid cacert")
			os.Exit(1)
		}
		transportStrategy = grpc.WithTransportCredentials(creds)
	}
	opts := []grpc.DialOption{transportStrategy}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewPeerServiceClient(conn)

	Download(client, addr, filePath)
}

// func main() {
// 	flag.Parse()

// 	// Default to insecure connection
// 	transportStrategy := grpc.WithInsecure()
// 	// Establish Secure connection using provided public key
// 	if *certFile != "" {
// 		creds, err := credentials.NewClientTLSFromFile(*certFile, "")
// 		if err != nil {
// 			fmt.Println("missing or invalid cacert")
// 			os.Exit(1)
// 		}
// 		transportStrategy = grpc.WithTransportCredentials(creds)
// 	}

// 	opts := []grpc.DialOption{transportStrategy}

// 	conn, err := grpc.Dial(*targetAddr, opts...)
// 	if err != nil {
// 		log.Fatalf("fail to dial: %v", err)
// 	}
// 	defer conn.Close()

// 	client := pb.NewPeerServiceClient(conn)

// 	Ping(client)
// }
