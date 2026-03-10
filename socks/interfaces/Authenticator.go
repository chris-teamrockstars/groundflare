package interfaces

import "io"
import "groundflare/socks/types"

type Authenticator interface {
	Authenticate(reader io.Reader, writer io.Writer, userAddr string) (*types.AuthContext, error)
	GetCode() uint8
}
