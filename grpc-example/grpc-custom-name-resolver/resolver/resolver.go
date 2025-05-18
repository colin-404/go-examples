package custom_resolver

import (
	"google.golang.org/grpc/resolver"
)

const (
	myScheme      = "example"
	myServiceName = "my-custom-service:1234"
	backendAddr   = "127.0.0.1:50051"
	backendAddr2  = "127.0.0.1:50052"
)

type myResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

type myResolverBuilder struct{}

func (*myResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	r := &myResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			myServiceName: {backendAddr, backendAddr2},
		},
	}
	r.start()
	return r, nil
}

func (*myResolverBuilder) Scheme() string { return myScheme }

func (r *myResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (r *myResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *myResolver) Close()                                {}

func init() {
	resolver.Register(&myResolverBuilder{})
}
