package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	pb "main/protos"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.String("port", os.Getenv("PORT"), "The server port")
	// port = flag.String("port", "10000", "The server port")
	default_cluster = flag.String("default_cluster", "11", "Default cluster")
)

var ctx = context.Background()

type ClusrerizationApiServer struct {
	pb.UnimplementedClusterizationAPIServer
}

func (s *ClusrerizationApiServer) UnaryClasterization(ctx context.Context, req *pb.GRPCRequest) (*pb.GRPCResponse, error) {
	return &pb.GRPCResponse{
		Pid:     req.GetPid(),
		Sid:     req.GetSid(),
		Cluster: rand.Int31(),
	}, nil
}

func newServer() *ClusrerizationApiServer {
	return &ClusrerizationApiServer{}
}

func (s *ClusrerizationApiServer) StreamClasterization(stream pb.ClusterizationAPI_StreamClasterizationServer) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		var pid = req.GetPid()
		var sid = req.GetSid()

		var key = strconv.FormatInt(pid, 10) + ":" + strconv.FormatInt(sid, 10)
		cl, err := rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			cl = *default_cluster
		} else if err != nil {
			panic(err)
		}
		cluster, _ := strconv.ParseInt(cl, 10, 32)
		res := &pb.GRPCResponse{
			Pid:     pid,
			Sid:     sid,
			Cluster: int32(cluster),
		}
		err = stream.Send(res)

	}
	return nil
}

func main() {
	flag.Parse()

	log.Printf("server: starting on port %s", *port)

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	// if *tls {
	// 	if *certFile == "" {
	// 		*certFile = data.Path("x509/server_cert.pem")
	// 	}
	// 	if *keyFile == "" {
	// 		*keyFile = data.Path("x509/server_key.pem")
	// 	}
	// 	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	// 	if err != nil {
	// 		log.Fatalf("Failed to generate credentials %v", err)
	// 	}
	// 	opts = []grpc.ServerOption{grpc.Creds(creds)}
	// }

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterClusterizationAPIServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
