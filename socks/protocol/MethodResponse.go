package protocol

import "fmt"
import "io"

// +-----+--------+
// | VER | METHOD |
// +-----+--------+
// |  1  |     1  |
// +-----+--------+
type MethodResponse struct {
	Version byte
	Method  byte
}

func NewMethodResponse(version byte, method byte) *MethodResponse {

	return &MethodResponse{
		Version: version,
		Method:  method,
	}

}

func ParseMethodResponse(reader io.Reader) (*MethodResponse, error) {

	version := []byte{0x00}
	_, err0 := reader.Read(version)

	if err0 == nil {

		method  := []byte{0x00}
		_, err1 := reader.Read(method)

		if err1 == nil {

			return &MethodResponse{
				Version: version[0],
				Method:  method[0],
			}, nil

		} else {
			return nil, fmt.Errorf("Missing SOCKS method in MethodRequest")
		}

	} else {
		return nil, fmt.Errorf("Missing SOCKS version in MethodRequest")
	}

}

func (response *MethodResponse) Bytes() []byte {

	bytes := make([]byte, 0, 2)
	bytes = append(bytes, response.Version)
	bytes = append(bytes, response.Method)

	return bytes

}
