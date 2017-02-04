package restish

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	testpb "github.com/philips/hacks/golang-rpc-comparison/pbpayload"
	"google.golang.org/grpc/grpclog"
)

type testClient struct {
	addr string
	*http.Client
}

type testServer struct {
}

func (s testServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	if err := req.Body.Close(); err != nil {
		panic(err)
	}
	in := testpb.SimpleRequest{}
	err = json.Unmarshal(body, &in)
	if err != nil {
		grpclog.Fatalf("json unmarshal failed %v: %v", req, err)
	}

	sr := &testpb.SimpleResponse{
		Payload: testpb.NewPayload(in.ResponseType, int(in.ResponseSize)),
	}
	buf, err := json.Marshal(sr)
	if err != nil {
		grpclog.Fatalf("json marshal failed %v: %v", req, err)
	}
	io.Copy(w, bytes.NewBuffer(buf))
}

// StartServer starts a gRPC server serving a benchmark service on the given
// address, which may be something like "localhost:0". It returns its listen
// address and a function to stop the server.
func StartServer(addr string) (string, func()) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		grpclog.Fatal(err)
	}

	s := &http.Server{
		Handler: testServer{},
	}

	go s.Serve(l)
	return l.Addr().String(), func() {
		l.Close()
	}
}

// DoCall performs an RPC with given stub and request and response sizes.
func DoCall(hc testClient, reqSize, respSize int) {
	pl := testpb.NewPayload(testpb.PayloadType_COMPRESSABLE, reqSize)
	req := &testpb.SimpleRequest{
		ResponseType: pl.Type,
		ResponseSize: int32(respSize),
		Payload:      pl,
	}
	buf, err := json.Marshal(req)
	if err != nil {
		grpclog.Fatalf("json marshal failed %v: %v", req, err)
	}
	sr := &testpb.SimpleResponse{}
	resp, err := hc.Post("http://"+hc.addr+"/", "test/json", bytes.NewBuffer(buf))
	if err != nil {
		grpclog.Fatalf("post failed %v: %v", resp, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		grpclog.Fatalf("read failed %v: %v", resp, err)
	}
	err = json.Unmarshal(body, sr)
	if err != nil {
		grpclog.Fatalf("json unmarshal failed %v: %v", resp, err)
	}
}
