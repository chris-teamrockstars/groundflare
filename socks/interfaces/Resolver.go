package interfaces

import "context"
import "net"

type Resolver interface {
	Resolve(ctx context.Context, name string) (context.Context, net.IP, error)
}
