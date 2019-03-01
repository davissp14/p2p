package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	client "github.com/davissp14/p2p/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	server "github.com/davissp14/p2p/pkg/server"
	pb "github.com/davissp14/p2p/pkg/service"
)

// Subcommands
var (
	serverCommand   = flag.NewFlagSet("server", flag.ExitOnError)
	pingCommand     = flag.NewFlagSet("ping", flag.ExitOnError)
	downloadCommand = flag.NewFlagSet("download", flag.ExitOnError)
	listCommand     = flag.NewFlagSet("list", flag.ExitOnError)
)

// Server Subcommands
var (
	serverPort     = serverCommand.Int("port", 8080, "server port")
	serverTLS      = serverCommand.Bool("tls", false, "enable secure transport")
	serverKeyFile  = serverCommand.String("key-file", "", "idenfity HTTPS client using this SSL key file")
	serverCertFile = serverCommand.String("cert-file", "", "identify HTTPS client using this SSL certificate file")
	serverHostname = serverCommand.String("hostname", "localhost", "identify HTTPS client using this SSL certificate file")
)

var (
	pingAddr     = pingCommand.String("addr", "", "ip address")
	pingCertFile = pingCommand.String("cert-file", "", "identify HTTPS client using this SSL certificate file")
)

// Download Subcommands
var (
	downloadAddr     = downloadCommand.String("addr", "", "ip address")
	downloadFilePath = downloadCommand.String("file-path", "", "ip address")
	downloadCertFile = downloadCommand.String("cert-file", "", "identify HTTPS client using this SSL certificate file")
)

// List Files Subcommands
var (
	listAddr      = listCommand.String("addr", "", "ip address")
	listDirectory = listCommand.String("dir", "/", "dir")
	listCertFile  = listCommand.String("cert-file", "", "identify HTTPS client using this SSL certificate file")
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf(
			`Your Network
Available Commands
==============
server:     Starts server.
ping        Ping Remote / Local node within your network.
download:   Exchange public keys with node    

`)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		serverCommand.Parse(os.Args[2:])
	case "ping":
		pingCommand.Parse(os.Args[2:])
	case "download":
		downloadCommand.Parse(os.Args[2:])
	case "list":
		listCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if serverCommand.Parsed() {
		var srv *grpc.Server
		if *serverTLS {
			fmt.Printf("Establishing secure connection on %s:%d\n", *serverHostname, *serverPort)
			creds, err := credentials.NewServerTLSFromFile(*serverCertFile, *serverKeyFile)
			if err != nil {
				fmt.Println(err.Error())
			}
			srv = grpc.NewServer(grpc.Creds(creds))
		} else {
			fmt.Printf("Establishing insecure connection on port %d\n", *serverPort)
			srv = grpc.NewServer()
		}
		pb.RegisterPeerServiceServer(srv, server.NewServer(*serverPort, *serverTLS, *serverKeyFile, *serverCertFile))

		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *serverHostname, *serverPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		srv.Serve(lis)
	}

	if pingCommand.Parsed() {
		conn, err := client.NewClientConn(*pingCertFile, *pingAddr)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		defer conn.Close()

		cl := pb.NewPeerServiceClient(conn)
		client.Ping(cl)
	}

	if downloadCommand.Parsed() {
		conn, err := client.NewClientConn(*downloadCertFile, *downloadAddr)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		defer conn.Close()

		cl := pb.NewPeerServiceClient(conn)
		client.Download(cl, *downloadAddr, *downloadFilePath)
	}

	if listCommand.Parsed() {
		conn, err := client.NewClientConn(*listCertFile, *listAddr)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		cl := pb.NewPeerServiceClient(conn)
		client.List(cl, *listDirectory)

	}

}
