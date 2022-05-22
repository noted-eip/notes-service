package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ForwardAuMetadatathUnaryInterceptor may be invoked on every incoming
// rpc before the actual endpoint code is reached and it shall forward
// the authorization header so calls to other services inherit the
// authorization state of the caller.
func ForwardAuthMetadatathUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return handler(IncomingToOutgoingContext(ctx), req)
}

// IncomingToOutgoingContext creates a new outgoing context with the metadata
// stored inside the incoming context.
func IncomingToOutgoingContext(parent context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(parent)
	return metadata.NewOutgoingContext(parent, md)
}
