package resolvers

import "context"
import "net"

type DNS struct {
}

func NewDNS() *DNS {
	return &DNS{}
}

func (resolver DNS) Resolve(ctx context.Context, hostname string) (context.Context, net.IP, error) {

	addr, err0 := net.ResolveIPAddr("ip", hostname)

	if err0 == nil {
		return ctx, addr.IP, nil
	} else {
		return ctx, nil, err0
	}

}
