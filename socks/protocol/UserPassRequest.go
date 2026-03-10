package protocol

import "io"
import "groundflare/socks/errors"

// +--------------+------+----------+------+----------+
// | USERPASS_VER | ULEN |   USER   | PLEN |   PASS   |
// +--------------+------+----------+------+----------+
// |      1       |   1  | Variable |   1  | Variable |
// +--------------+------+----------+------+----------+
type UserPassRequest struct {
	Version  byte
	Username string // 1-255 bytes
	Password string // 1-255 bytes
}

func NewUserPassRequest(version byte, username string, password string) UserPassRequest {

	return UserPassRequest{
		Version:  version,
		Username: username,
		Password: password,
	}

}

func ParseUserPassRequest(reader io.Reader) (*UserPassRequest, error) {

	version         := []byte{0x00}
	username_length := []byte{0x00}

	// Get the version and username length
	_, err0 := reader.Read(version)
	_, err1 := reader.Read(username_length)

	if err0 == nil && err1 == nil {

		if version[0] == AuthVersion {

			username := make([]byte, int(username_length[0]))
			_, err2 := io.ReadAtLeast(reader, username, int(username_length[0]))

			if err2 == nil {

				password_length := []byte{0x00}
				_, err3 := reader.Read(password_length)

				if err3 == nil && int(password_length[0]) > 0 {

					password := make([]byte, int(password_length[0]))
					_, err4 := io.ReadAtLeast(reader, password, int(password_length[0]))

					if err4 == nil {

						return &UserPassRequest{
							Version:  version[0],
							Username: string(username),
							Password: string(password),
						}, nil

					} else {
						return nil, err4
					}

				} else {
					return nil, err3
				}

			} else {
				return nil, err2
			}

		} else {
			return nil, errors.AuthUnsupportedVersion
		}

	} else {
		return nil, err0
	}

}

func (request *UserPassRequest) Bytes() []byte {

	bytes := make([]byte, 0, 2 + len(request.Username) + 1 + len(request.Password))
	bytes = append(bytes, request.Version)
	bytes = append(bytes, byte(len(request.Username)))
	bytes = append(bytes, []byte(request.Username)...)
	bytes = append(bytes, byte(len(request.Password)))
	bytes = append(bytes, []byte(request.Password)...)

	return bytes

}

