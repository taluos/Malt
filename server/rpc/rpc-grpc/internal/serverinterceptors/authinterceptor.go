// auth check by jwt
package serverinterceptors

import (
	"context"

	"github.com/taluos/Malt/pkg/auth-jwt"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"
	"github.com/taluos/Malt/pkg/log"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
)

func SteamAuthorizeInterceptor(keyFunc jwt.Keyfunc, authenticator *auth.Authenticator) grpc.StreamServerInterceptor {
	if authenticator == nil {
		var err error
		authenticator, err = auth.NewAuthenticator(keyFunc)
		if err != nil {
			log.Errorf("failed to create authenticator: %s", err)
			return nil
		}
	}
	return func(svr any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		if err := authenticator.RPCAuthenticate(stream.Context(), "", info.FullMethod, ""); err != nil {
			log.Errorf("auth failed: %s", err)
			return errors.WithCode(code.ErrInvalidAuthHeader, "auth failed")
		}
		return handler(svr, stream)
	}
}

func UnaryAuthorizeInterceptor(keyFunc jwt.Keyfunc, authenticator *auth.Authenticator) grpc.UnaryServerInterceptor {
	if authenticator == nil {
		var err error
		authenticator, err = auth.NewAuthenticator(keyFunc)
		if err != nil {
			log.Errorf("failed to create authenticator: %s", err)
			return nil
		}
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if err := authenticator.RPCAuthenticate(ctx, "", info.FullMethod, ""); err != nil {
			log.Errorf("auth failed: %s", err)
			return nil, errors.WithCode(code.ErrInvalidAuthHeader, "auth failed")
		}
		return handler(ctx, req)
	}
}
