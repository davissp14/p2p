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
	downloadCommand = flag.NewFlagSet("download", flag.ExitOnError)
)

// Server Subcommands
var (
	serverPort     = serverCommand.Int("port", 8080, "server port")
	serverTLS      = serverCommand.Bool("tls", false, "enable secure transport")
	serverKeyFile  = serverCommand.String("key-file", "", "idenfity HTTPS client using this SSL key file")
	serverCertFile = serverCommand.String("cert-file", "", "identify HTTPS client using this SSL certificate file")
)

var (
	downloadAddr     = downloadCommand.String("addr", "", "ip address")
	downloadFilePath = downloadCommand.String("file-path", "", "ip address")
	downloadCertFile = downloadCommand.String("cert-file", "", "identify HTTPS client using this SSL certificate file")
)

// addNodeCommand := flag.NewFlagSet("add_node", flag.ExitOnError)
// listNodesCommand := flag.NewFlagSet("list_nodes", flag.ExitOnError)
// pingCommand := flag.NewFlagSet("ping", flag.ExitOnError)
// listFilesCommand := flag.NewFlagSet("list_files", flag.ExitOnError)
// downloadCommand := flag.NewFlagSet("download", flag.ExitOnError)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf(
			`Your Network
Available Commands
==============
server:     Starts server.
ping        Ping Remote / Local node within your network.
download:   Exchange public keys with node    

Optional Flags
cacert:  Required for most commands when tls is configured
`)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		serverCommand.Parse(os.Args[2:])
	case "download":
		downloadCommand.Parse(os.Args[2:])

	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if serverCommand.Parsed() {
		var srv *grpc.Server
		if *serverTLS {
			fmt.Printf("Establishing secure connection on port %d\n", *serverPort)
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

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *serverPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		srv.Serve(lis)
	}

	if downloadCommand.Parsed() {
		client.NewClient(*downloadCertFile, *downloadAddr, *downloadFilePath)
	}

}
