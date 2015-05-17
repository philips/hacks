package main

import (
	"flag"
	"net/http"

	"github.com/gengo/grpc-gateway/runtime"
	"github.com/golang/glog"
	"golang.org/x/net/context"

	gw "github.com/philips/hacks/grpc-play/gateway"
)

var (
	echoEndpoint = flag.String("echo_endpoint", "localhost:10000", "endpoint of YourService")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	err := gw.RegisterYourServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint)
	if err != nil {
		return err
	}

	http.ListenAndServe(":8080", mux)
	return nil
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
