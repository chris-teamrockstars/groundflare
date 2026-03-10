package socks

import "testing"
import "github.com/stretchr/testify/assert"

func TestStaticCredentials(t *testing.T) {
	creds := StaticCredentials{
		"foo": "bar",
		"baz": "",
	}

	assert.True(t, creds.Valid("foo", "bar", ""))
	assert.True(t, creds.Valid("baz", "", ""))
	assert.False(t, creds.Valid("foo", "", ""))
}
