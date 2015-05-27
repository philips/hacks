package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/philips/hacks/grpc-play/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "testdata/ca.pem", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		var sn string
		if *serverHostOverride != "" {
			sn = *serverHostOverride
		}
		var creds credentials.TransportAuthenticator
		if *caFile != "" {
			var err error
			creds, err = credentials.NewClientTLSFromFile(*caFile, sn)
			if err != nil {
				grpclog.Fatalf("Failed to create TLS credentials %v", err)
			}
		} else {
			creds = credentials.NewClientTLSFromCert(nil, sn)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	opts = append(opts, grpc.WithDialer(httpDial))
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewYourServiceClient(conn)

	msg, err := client.Echo(context.Background(), &pb.StringMessage{os.Args[1]})
	println(msg.Value)
}

func httpDial(addr string, timeout time.Duration) (net.Conn, error) {
	log.Println("dial...")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("write HTTP line")
	_, err = conn.Write([]byte("GET /grpc HTTP/1.1\r\n"))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("write Header line")
	_, err = conn.Write([]byte("HOST: " + addr + "\r\n\r\n"))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// handshake
	b := make([]byte, 2)
	_, err = conn.Read(b)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
