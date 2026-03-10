package socks

import "context"
import "testing"
import "github.com/stretchr/testify/assert"
import "github.com/stretchr/testify/require"

func TestDNSResolver(t *testing.T) {
	d := DNSResolver{}
	ctx := context.Background()

	_, addr, err := d.Resolve(ctx, "localhost")
	require.NoError(t, err)
	assert.True(t, addr.IsLoopback())
}
