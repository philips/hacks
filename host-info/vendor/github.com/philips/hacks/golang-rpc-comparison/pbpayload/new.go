package pbpayload

import (
	"google.golang.org/grpc/grpclog"
)

func NewPayload(t PayloadType, size int) *Payload {
	if size < 0 {
		grpclog.Fatalf("Requested a response with invalid length %d", size)
	}
	body := make([]byte, size)
	switch t {
	case PayloadType_COMPRESSABLE:
	case PayloadType_UNCOMPRESSABLE:
		grpclog.Fatalf("PayloadType UNCOMPRESSABLE is not supported")
	default:
		grpclog.Fatalf("Unsupported payload type: %d", t)
	}
	return &Payload{
		Type: t,
		Body: body,
	}
}
