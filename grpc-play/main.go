package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/gengo/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	gw "github.com/philips/hacks/grpc-play/gateway"
	pb "github.com/philips/hacks/grpc-play/proto"
)

var (
	tls       = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile  = flag.String("cert_file", "testdata/server1.pem", "The TLS cert file")
	keyFile   = flag.String("key_file", "testdata/server1.key", "The TLS key file")
	port      = flag.Int("port", 10000, "The rpc server port")
	proxyPort = flag.Int("proxy_port", 8080, "The proxy server port")
)

type myService struct{}

func (m *myService) Echo(c context.Context, s *pb.StringMessage) (*pb.StringMessage, error) {
	fmt.Printf("rpc request Echo(%q)\n", s.Value)
	return s, nil
}

func newServer() *myService {
	return new(myService)
}

func proxy() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	err := gw.RegisterYourServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		fmt.Printf("proxy: %v\n", err)
		return
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", *proxyPort), mux)
	fmt.Printf("proxy: %v\n", err)
	return
}

func serve(opts []grpc.ServerOption) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterYourServiceServer(grpcServer, newServer())

	err = grpcServer.Serve(lis)
	fmt.Printf("proxy: %v\n", err)
	return
}

func main() {
	flag.Parse()
	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	go serve(opts)
	fmt.Printf("grpc on port: %d\n", *port)
	fmt.Printf("rest on port: %d\n", *proxyPort)
	proxy()
}
