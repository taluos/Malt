package metadata

import (
	"context"
	"fmt"
)

type clientMetadataKey struct{}

// NewClientContext creates a new context with client md attached.
// This function is implemented in the same way as in NewServerContext and grpc.NewIncomingContext
func NewClientContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, clientMetadataKey{}, md)
}

// FromClientContext returns the client metadata in ctx if it exists.
// 从context中获取client metadata
func FromClientContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(clientMetadataKey{}).(Metadata)
	return md, ok
}

// AppendToClientContext returns a new context with the provided kv merged
// with any existing metadata in the context.
// This functions is a modification of grpc.AppendToOutgoingContext
func AppendToClientContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendToClientContext got an odd number of input pairs for metadata: %d", len(kv)))
	}
	md, _ := FromClientContext(ctx)
	md = md.DeepClone()
	for i := 0; i < len(kv); i += 2 {
		md.Set(kv[i], kv[i+1])
	}
	return NewClientContext(ctx, md)
}

// MergeToClientContext merge new metadata into ctx.
func MergeToClientContext(ctx context.Context, cmd Metadata) context.Context {
	md, _ := FromClientContext(ctx)
	//md = md.Clone()
	//for k, v := range cmd {
	//	md[k] = v
	//}
	md = md.DeepClone()
	return NewClientContext(ctx, md)
}
