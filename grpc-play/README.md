## gRPC + REST Gateway Play

This runs a simple gRPC service on port 10000 and a REST interface to this gRPC endpoint on port 8080.

To try it all out do this:

```
$ go install github.com/philips/hacks/grpc-play/grpc-play-client
$ go install github.com/philips/hacks/grpc-play

$ grpc-play
$ grpc-play-client "my first rpc echo"
$ curl -X POST localhost:8080/v1/example/echo -H "Content-Type: text/plain"  -d '{"value":"my REST echo"}'
{"value":"my REST echo"}
```


Huge thanks to the hard work people have put into the [Go gRPC bindings][gogrpc] and [gRPC to JSON Gateway][grpcgateway]

[gogrpc]: https://github.com/grpc/grpc-go
[grpcgateway]: https://github.com/gengo/grpc-gateway
