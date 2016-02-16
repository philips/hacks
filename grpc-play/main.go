package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gengo/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "github.com/philips/hacks/grpc-play/proto"
)

var (
	tls      = flag.Bool("tls", true, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "certs/server.pem", "The TLS cert file")
	keyFile  = flag.String("key_file", "certs/server.key", "The TLS key file")
	port     = flag.Int("port", 10000, "The rpc server port")
)

type myService struct{}

func (m *myService) Echo(c context.Context, s *pb.StringMessage) (*pb.StringMessage, error) {
	fmt.Printf("rpc request Echo(%q)\n", s.Value)
	return s, nil
}

func newServer() *myService {
	return new(myService)
}

// GRPCHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func serve(opts []grpc.ServerOption) {
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterYourServiceServer(grpcServer, newServer())
	ctx := context.Background()

	name := fmt.Sprintf("localhost:%d", *port)
	dcreds, err := credentials.NewClientTLSFromFile(*certFile, name)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return
	}
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	mux := runtime.NewServeMux()
	err = pb.RegisterYourServiceHandlerFromEndpoint(ctx, mux, name, dopts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return
	}

	err = http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), *certFile, *keyFile, grpcHandlerFunc(grpcServer, mux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return
}

func main() {
	flag.Parse()
	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	fmt.Printf("grpc on port: %d\n", *port)
	serve(opts)
}
