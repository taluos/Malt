package grpc

import (
	"context"
	"testing"
	"time"

	// rpcserver "github.com/taluos/Malt/server/rpc/rpc-grpc"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	s := NewServer(WithAddress(defaultAddress))
	assert.NotNil(t, s, "Expected non-nil server")
}

func TestServerCreation(t *testing.T) {
	var err error
	s := NewServer(
		WithAddress(defaultAddress),
		// rpcserver.WithListener(listener),
	)

	if s == nil {
		t.Fatal("Expected non-nil server")
	}
	var ctx = context.Background()
	go func() {
		err = s.Start(ctx)
		assert.NoError(t, err, "Expected no error when starting server")
	}()
	time.Sleep(10 * time.Second)

	s.Stop(ctx)
}
