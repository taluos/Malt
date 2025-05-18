# How to build a interception ?

## Interceptor Types

In gRPC, there are two main types of interceptors:

1. UnaryInterceptor
2. StreamInterceptor

## Unary Interceptor Example

### Example1 Input Requirements

- `ctx context.Context`: The context containing metadata and deadlines
- `req any`: The request message from the client (type will be determined by your proto definition)
- `info *grpc.UnaryServerInfo`: Contains metadata about the RPC call (method name, server instance)
- `handler grpc.UnaryHandler`: The actual RPC method implementation to be called

### Example1 Output Requirements

- Returns `(any, error)`:
  - First return value: The response message to be sent back to the client
  - Second return value: Any error that occurred during processing

```go
func UnaryExampleInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
    // Pre-processing
    log.Printf("Before handling request: %v", info.FullMethod)
    
    // Call the actual handler
    resp, err := handler(ctx, req)
    
    // Post-processing
    if err != nil {
        log.Printf("Error occurred: %v", err)
    }
    
    return resp, err
}
```

## Stream Interceptor Example

### Example2 Input Requirements

- `srv any`: The server implementation
- `ss grpc.ServerStream`: The server-side stream interface for sending/receiving messages
- `info *grpc.StreamServerInfo`: Contains metadata about the stream RPC call
- `handler grpc.StreamHandler`: The actual stream RPC method implementation

### Example2 Output Requirements

- Returns `error`: Any error that occurred during stream processing

```go
func StreamExampleInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    // Pre-processing
    log.Printf("Before handling stream: %v", info.FullMethod)
    
    // Call the actual handler
    err := handler(srv, ss)
    
    // Post-processing
    if err != nil {
        log.Printf("Stream error occurred: %v", err)
    }
    
    return err
}
```

## Registering Interceptors

Register interceptors when creating the gRPC server:

```go
server := grpc.NewServer(
    grpc.UnaryInterceptor(UnaryExampleInterceptor),
    grpc.StreamInterceptor(StreamExampleInterceptor),
)
```

## Common Use Cases

1. Authentication and Authorization
2. Logging
3. Error Handling
4. Timeout Control
5. Rate Limiting
6. Metrics Collection

## Best Practices

1. Keep interceptors simple and focused
2. Avoid time-consuming operations in interceptors
3. Handle errors and exceptions properly
4. Pay attention to context propagation
5. Consider interceptor execution order
