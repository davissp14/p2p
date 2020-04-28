package server

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"fmt"
	"path/filepath"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes"

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
		return fmt.Errorf(err.Error())
	}

	for _, f := range fs {
		path_to_file := filepath.Clean(fmt.Sprintf("%s/%s", req.Directory, f.Name()))
		fi, err := os.Lstat(path_to_file)
		if err != nil {
			fmt.Println(fmt.Sprintf("%s : %s", "error", err.Error()))
		}

		mode := f.Mode()
		time, _ := ptypes.TimestampProto(f.ModTime())

		file := pb.File{
			Name:  f.Name(),
			Filepath: path_to_file,
			ModTime: time, 
			Mode: mode.String(),
			Symlink: (fi.Mode()&os.ModeSymlink != 0),
			LinkedTo: "",
			ValidLink: false,
			IsDir: false,
		}


		// file := pb.File{
		// 	Filepath: req.Directory,
		// 	ModTime: time,
		// 	Mode: mode.String(),
		// 	IsDir: mode.IsDir(),
		// 	Size:  f.Size(),
		// 	IsRegular: mode.IsRegular(),
		// }

		// True if the file is a symlink.
		if file.Symlink {
			origin, err := os.Readlink(path_to_file)
			if err != nil {
		 		fmt.Println(fmt.Sprintf("%s : %s", "error", err.Error()))
			}
			path_to_origin := filepath.Clean(fmt.Sprintf("/%s", origin))
		
			// Set LinkedTo
			file.LinkedTo = path_to_origin 
			linked, err := os.Lstat(path_to_origin)
			if err == nil { 
				file.ValidLink = true  
				file.IsDir = linked.IsDir()
			}
		} else {
			file.IsDir = f.IsDir() 
		}

		if !file.IsDir {
			file.Size = f.Size()
		}

		stream.Send(&file)
	}

	return nil
}
