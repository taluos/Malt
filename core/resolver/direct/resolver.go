package direct

import "google.golang.org/grpc/resolver"

type directResolver struct {
	// addresses []resolver.Address
	cc resolver.ClientConn
}

func NewDirectResolver() *directResolver {
	return &directResolver{}
}

func (r *directResolver) Close() {

}

func (r *directResolver) ResolveNow(opt resolver.ResolveNowOptions) {

}
