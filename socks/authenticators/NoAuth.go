package authenticators

import "io"
import "groundflare/socks/protocol"
import "groundflare/socks/types"

type NoAuth struct{
}

func (authenticator NoAuth) GetCode() uint8 {
	return protocol.MethodNoAuth
}

func (authenticator NoAuth) Authenticate(_ io.Reader, writer io.Writer, _ string) (*types.AuthContext, error) {

	_, err := writer.Write([]byte{protocol.Version5, protocol.MethodNoAuth})

	return &types.AuthContext{
		Method:  protocol.MethodNoAuth,
		Payload: make(map[string]string),
	}, err

}

