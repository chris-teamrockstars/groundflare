package authenticators

import "io"
import "groundflare/socks/errors"
import "groundflare/socks/interfaces"
import "groundflare/socks/protocol"
import "groundflare/socks/types"

type UserPass struct {
	Credentials interfaces.Credentials
}

func (authenticator UserPass) GetCode() uint8 {
	return protocol.MethodUserPassAuth
}

func (authenticator UserPass) Authenticate(reader io.Reader, writer io.Writer, userAddr string) (*types.AuthContext, error) {

	if _, err := writer.Write([]byte{protocol.Version5, protocol.MethodUserPassAuth}); err != nil {
		return nil, err
	}

	request, err := protocol.ParseUserPassRequest(reader)

	if err != nil {
		return nil, err
	}

	if !authenticator.Credentials.Valid(request.Username, request.Password, userAddr) {

		if _, err := writer.Write([]byte{protocol.AuthVersion, protocol.AuthFailure}); err != nil {
			return nil, err
		}

		return nil, errors.AuthFailure

	}

	if _, err := writer.Write([]byte{protocol.AuthVersion, protocol.AuthSuccess}); err != nil {
		return nil, err
	}

	return &types.AuthContext{
		protocol.MethodUserPassAuth,
		map[string]string{
			"username": request.Username,
			"password": request.Password,
		},
	}, nil

}

