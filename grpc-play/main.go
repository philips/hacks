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

const (
	port = 10000
	key  = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4akInId7aN716BqD0aqIxRcGiVADCy/P+U60YNGE1k/lzfok
a+U15Tsj3P3QWHDnhwYk8rRn8XcbrSqa0UreeD56ZlS7aW8OMDfMu4du1tp9dDJQ
bDGeFQeJnQaqXZRvb9ZhvWneXukIqG2puTb4czjiE3Ns2LBn4EtOVYj22rZ7Ycxg
CWKlE3cd4c4aeLcutMgw2oVsAT7kVRH1a7sxZLxxv4Z8O73Y3pF8hnoopeXq41dh
wq/8pRf4sYGGWaUWzyYqaTVCAIvSrDcQTRNRfoi1OlAwQfkPtIEQc+pXmrnOXoHq
FGSM2E/moRhcsaITMW7IoFAEXFgTLVNyJFKDawIDAQABAoIBAGa6t/PibZkZX2NR
ps7tTqRCzjP2m9wc4RYC8XTeDRYve1Og27HOwBFiMfpqBc4tYAmFD475+BPiAFR9
/8rlxY+CCeDKFCN0bkYfMPHyPtBK2S/rs9b9Y5T5QHmezjIZ3/1O4GIFbzlP5yQe
AFUUJFj3/VuTgWrILBHc1oVz+8Um0OP5ykMk5Kbl/n5PVyVp0v5SNS9I1T2EGfNc
o16uTwj+Qo4Pe3B7VCq/3lnx9lc819l1Q7heC/l9467EuZn996b7sW2A5yrWBarL
jYwThJVu6JTUbLj/GXCS+ZE/6zCJ41i6OzTcNCsksc9SW9sdblPfIdmc5DFnLXqH
dC45adkCgYEA+A27o46lS7olVQR/dMfNNhHdMKxE2lNQm6fA0F+KoTZntiBRksEH
57NjZ/ZArNksyNI5uP88Ykhifr6k2b80xPlhfDw4ZDl8CRXiA771UiHsIvB0wD9N
AcKvUwvFOv341+pv4PUl0ohjuLCiVfJJNZ4+dSt8BjiUtXsaC9oXfJcCgYEA6OOn
lqPQM8JogKO7ukCLK4sdJkhWPSznCbWKZ6WcvYbGygRsGQHWvWbU3iOBr682ki+e
X8oNT8Rary9pSemZCY2YMgW80vfmgJGK2sYsVqbhYxs/acBqL8WXehUDPZvt6sUg
L3wbyaJQX+mCZeAse7YiIKugJGxpSeUU4iKthk0CgYEAxzZ4QKW5+LRZcQr4tbgV
BdyY8JMZhOGudiPmhTKF6m0AI9OWz655A8sdBYxOasLL5ch4FSveuev6NmIzkLCv
15WUhiry+wLzq3RInMuKx9h4haLpkNAFr2lEVwS39GWtqPIdweP/6TIiLFynMzEv
PIGHFaDDrVdZjtp7k5Mmk6cCgYAYNdptOZhiWRp+DQduBFmzbCHaofh9IZbfFoVN
4xSZS1KNG8qVCvDk/bSxZyWLOv7EUbj4Ikwh97qprZcXfPZQ3OxuftQzZlwLD5ZM
yf8//tc9c06zUrJ3RuZJZbfRhs1D87w104QcAQiz/9Vze8uEDNodZVofjzme2fbC
z3IUnQKBgGtSSSRh95AM4s+DwIFZurVJ2l4iXRooTlxORq8TKkhCZNTCU6PA+hF6
NHZRXR7Jp+mfW7Q58DDv9LsEu/JLZhw7BP8tdWdmn82dRyG1Dwc+53CrN2RpPAdo
vSlA6AxOSTJgZL5/A0g862lI5qs6cq2seuB4qnCV0h5dPP/7lAbN
-----END RSA PRIVATE KEY-----
`
	cert = `-----BEGIN CERTIFICATE-----
MIIEBjCCAu6gAwIBAgIJAKaeXTcX/kNiMA0GCSqGSIb3DQEBBQUAMF8xCzAJBgNV
BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQxGDAWBgNVBAMTD2xvY2FsaG9zdDoxMDAwMDAeFw0xNjAy
MTUyMjE3NTlaFw0yNjAyMTIyMjE3NTlaMF8xCzAJBgNVBAYTAkFVMRMwEQYDVQQI
EwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQx
GDAWBgNVBAMTD2xvY2FsaG9zdDoxMDAwMDCCASIwDQYJKoZIhvcNAQEBBQADggEP
ADCCAQoCggEBAOGpCJyHe2je9egag9GqiMUXBolQAwsvz/lOtGDRhNZP5c36JGvl
NeU7I9z90Fhw54cGJPK0Z/F3G60qmtFK3ng+emZUu2lvDjA3zLuHbtbafXQyUGwx
nhUHiZ0Gql2Ub2/WYb1p3l7pCKhtqbk2+HM44hNzbNiwZ+BLTlWI9tq2e2HMYAli
pRN3HeHOGni3LrTIMNqFbAE+5FUR9Wu7MWS8cb+GfDu92N6RfIZ6KKXl6uNXYcKv
/KUX+LGBhlmlFs8mKmk1QgCL0qw3EE0TUX6ItTpQMEH5D7SBEHPqV5q5zl6B6hRk
jNhP5qEYXLGiEzFuyKBQBFxYEy1TciRSg2sCAwEAAaOBxDCBwTAdBgNVHQ4EFgQU
gFVW7OdpoctFELBHiiER39867owwgZEGA1UdIwSBiTCBhoAUgFVW7OdpoctFELBH
iiER39867oyhY6RhMF8xCzAJBgNVBAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRl
MSEwHwYDVQQKExhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxGDAWBgNVBAMTD2xv
Y2FsaG9zdDoxMDAwMIIJAKaeXTcX/kNiMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcN
AQEFBQADggEBAIG1Pq8zLTDrgaawXPbpJhs8LybuYppYr51lBR1Jc4mI8XDshE5J
wcpqAhWm3jZMTokDNSA50hMNpUPhUmk8JQIgjz5G9V7ETVWx1QbfdEr6deqOq6o2
itwt5ZmlXC0ZN+zyt0MS7NQnCHM9jpb5MpGCDWUyCdbArNGp+Mj78P+rC6P/i02J
0RdF6WCe/VyDnw8OgRndQDD4U4jfO2isWN8OnXh8fWAgXJJTtuZLtZPZXAbFtQNn
lzCMcbPnvl2Yj1naUDH7u/A4XYV/1dIetMI8E0Ef8UNwtlEtUW1I9wIrM1/eQ9d5
gyRJacQ9aPzPnC8wdNWeBw/7QNmrZbseRgE=
-----END CERTIFICATE-----
`
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

	name := fmt.Sprintf("localhost:%d", port)
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

	config := &tls.Config{
		MinVersion: tls.VersionTLS10,
	}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return err
	}

	conn, err := net.Listen("tcp", name)
	if err != nil {
		return err
	}

	listener := tls.NewListener(conn, config)

	err = http.Serve(listener, grpcHandlerFunc(grpcServer, mux))
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

	fmt.Printf("grpc on port: %d\n", port)
	serve(opts)
}
