package main

import (
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "github.com/philips/hacks/grpc-play/proto"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "testdata/server1.pem", "The TLS cert file")
	keyFile  = flag.String("key_file", "testdata/server1.key", "The TLS key file")
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

func proxy(gopts []grpc.ServerOption) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	/*
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}
		err := pb.RegisterYourServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", *port), opts)
		if err != nil {
			fmt.Printf("proxy: %v\n", err)
			return
		}
	*/

	grpcServer := grpc.NewServer(gopts...)
	pb.RegisterYourServiceServer(grpcServer, newServer())

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("request")
		if r.Header.Get("content-type") != "application/grpc" {
			grpcServer.ServeHTTP(w, r)
			return
		}
		grpcServer.ServeHTTP(w, r)
		//		mux.ServeHTTP(w, r)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", *port), Handler: handler,
		TLSConfig: &tls.Config{
			CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		},
	}
	http2.ConfigureServer(server, nil)
	err := server.ListenAndServe()
	fmt.Printf("proxy: %v\n", err)
	return
}

func main() {
	flag.Parse()
	var gopts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		gopts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	fmt.Printf("port: %d\n", *port)
	proxy(gopts)
}
