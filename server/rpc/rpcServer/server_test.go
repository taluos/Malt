package rpcserver_test

import (
	"context"
	"net"
	"testing"
	"time"

	rpcserver "github.com/taluos/Malt/server/rpc/rpcServer"

	"github.com/stretchr/testify/assert"
)

var add = "127.0.0.1:8081"

var listener net.Listener

func TestNewServer(t *testing.T) {
	s := rpcserver.NewServer(rpcserver.WithAddress(add))
	assert.NotNil(t, s, "Expected non-nil server")
}

func TestServerCreation(t *testing.T) {
	var err error
	s := rpcserver.NewServer(
		rpcserver.WithAddress(add),
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
	err = s.Stop(ctx)
	assert.NoError(t, err, "Expected no error when stopping server")
}
