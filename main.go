package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"

	pb "github.com/Nataliyi/go_grpc/protos"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", os.Getenv("PORT"), "The server port")
)

type ClusrerizationApiServer struct {
	pb.UnimplementedClusterizationAPIServer
}

func (s *ClusrerizationApiServer) UnaryClasterization(ctx context.Context, req *pb.GRPCRequest) (*pb.GRPCResponse, error) {
	return &pb.GRPCResponse{
		Pid:     req.GetPid(),
		Sid:     req.GetSid(),
		Cluster: rand.Intn(100),
	}, nil
}

func (s *ClusrerizationApiServer) StreamClasterization(stream pb.ClusterizationAPI_StreamClasterizationServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		res := &pb.GRPCResponse{
			Pid:     req.GetPid(),
			Sid:     req.GetSid(),
			Cluster: rand.Intn(100),
		}
		err = stream.Send(res)
	}
	return nil
}

func newServer() *ClusrerizationApiServer {
	return &ClusrerizationApiServer{}
}

func main() {
	flag.Parse()

	log.Printf("server: starting on port %s", *port)
	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	if *tls {
		if *certFile == "" {
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = data.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterClusterizationAPIServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
