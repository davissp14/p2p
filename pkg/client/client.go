package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/davissp14/p2p/pkg/service"

	"google.golang.org/grpc"
)

var (
	targetAddr = flag.String("targetAddr", "localhost:8080", "endpoint of your service")
)

func Ping(client pb.PeerServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ping := pb.PingMessage{Message: "PING"}

	resp, err := client.Ping(ctx, &ping)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resp.GetMessage())
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*targetAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewPeerServiceClient(conn)

	Ping(client)

}
