package metadata

import (
	"context"
)

type serverMetadataKey struct{}

// NewServerContext creates a new context with client md attached.
// This function is implemented in the same way as NewClientContext and grpc.NewIncomingContext
func NewServerContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, serverMetadataKey{}, md)
}

// FromServerContext returns the server metadata in ctx if it exists.
// This function is a modification of grpc.fromOutgoingContextRaw
func FromServerContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(serverMetadataKey{}).(Metadata)
	return md, ok
}
