package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/davissp14/p2p/pkg/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	// Exclude the filepath from the name
	name := filepath.Base(filePath)
	file, err := os.Create(fmt.Sprintf("downloads/%s", name))
	if err != nil {
		log.Fatalf("Failed to create cert file: error. %s", err.Error())
		return
	}

	defer file.Close()

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		_, err = file.Write(in.Data)
		if err != nil {
			log.Fatalf("Failed to stream file. error: %s", err.Error())
			return
		}
	}
	log.Printf("File `%s` download was a success!\n", name)
}

func List(client pb.PeerServiceClient, directory string) ([]*pb.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := pb.ListRequest{
		Directory: directory,
	}
	stream, err := client.List(ctx, &req)
	if err != nil {
		return nil, err
	}

	var files []*pb.File

	for {
		file, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
			return nil, err
		}
		files = append(files, file)
	}


	for _, file := range files {
		if file.IsDir {
			fmt.Printf("%-60s | %-10s | %-10s | %-10s\n", file.Name, "", strconv.FormatBool(file.IsDir), file.LinkedTo)
		} else {
			fmt.Printf("%-60s | %-10s  | %-10s | %-10s\n", file.Name, fmt.Sprintf("%d bytes", file.Size), strconv.FormatBool(file.IsDir), file.LinkedTo)
		}
	}

	return files, nil 

}

func NewClientConn(cert, addr string) (*grpc.ClientConn, error) {
	// Default to insecure connection
	transportStrategy := grpc.WithInsecure()
	// Establish Secure connection using provided public key
	if cert != "" {
		fmt.Println("Secure client conn")
		creds, err := credentials.NewClientTLSFromFile(cert, "")
		if err != nil {
			return nil, err
		}
		transportStrategy = grpc.WithTransportCredentials(creds)
	}
	opts := []grpc.DialOption{transportStrategy}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
