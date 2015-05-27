package main

import (
	"flag"
	"fmt"
	"log"
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
	mux := runtime.NewServeMux()
	err := gw.RegisterYourServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		fmt.Printf("proxy: %v\n", err)
		return
	}

	http.Handle("/", mux)
	return
}

func serve(opts []grpc.ServerOption, ch chan net.Conn) {
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterYourServiceServer(grpcServer, newServer())

	grpcServer.Serve(&grpcListener{ch: ch})
	return
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	lch := make(chan net.Conn, 128)

	http.HandleFunc("/grpc", func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}
		conn, rwbuf, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// handshake
		rwbuf.Write([]byte("OK"))
		rwbuf.Flush()
		log.Println("accept grpc")
		lch <- conn
	})

	go http.Serve(lis, nil)

	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	fmt.Printf("port: %d\n", *port)
	go serve(opts, lch)
	go proxy()
	select {}
}

type grpcListener struct {
	ch   chan net.Conn
	addr net.Addr
}

func (gl *grpcListener) Accept() (net.Conn, error) {
	return <-gl.ch, nil
}

func (gl *grpcListener) Close() error {
	close(gl.ch)
	return nil
}

func (gl *grpcListener) Addr() net.Addr {
	return gl.addr
}
