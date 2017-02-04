package main

import (
	"flag"
	"os"

	pb "github.com/philips/hacks/grpc-play/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

var (
	tls    = flag.Bool("tls", true, "Connection uses TLS if true, else plain TCP")
	caFile = flag.String("ca_file", "testdata/ca.pem", "The file containning the CA root cert file")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	var creds credentials.TransportAuthenticator
	if *caFile != "" {
		var err error
		creds, err = credentials.NewClientTLSFromFile(*caFile, sn)
		if err != nil {
			grpclog.Fatalf("Failed to create TLS credentials %v", err)
		}
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewYourServiceClient(conn)

	msg, err := client.Echo(context.Background(), &pb.StringMessage{os.Args[1]})
	println(msg.Value)
}
