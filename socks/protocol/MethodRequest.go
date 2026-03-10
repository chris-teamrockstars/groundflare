package protocol

import "fmt"
import "io"

// +-----+----------+---------------+
// | VER | NMETHODS |    METHODS    |
// +-----+----------+---------------+
// |  1  |     1    | X'00' - X'FF' |
// +-----+----------+---------------+
type MethodRequest struct {
	Version byte
	Methods []byte
}

func NewMethodRequest(version byte, methods []byte) *MethodRequest {

	return &MethodRequest{
		Version: version,
		Methods: methods,
	}

}

func ParseMethodRequest(reader io.Reader) (*MethodRequest, error) {

	version := []byte{0x00}
	_, err0 := reader.Read(version)

	if err0 == nil {

		methods_length := []byte{0x00}
		_, err1        := reader.Read(methods_length)

		if err1 == nil && int(methods_length[0]) > 0 {

			methods := make([]byte, int(methods_length[0]))
			_, err2 := io.ReadAtLeast(reader, methods, int(methods_length[0]))

			if err2 == nil {

				return &MethodRequest{
					Version: version[0],
					Methods: methods,
				}, nil

			} else {
				return nil, fmt.Errorf("Missing SOCKS methods in MethodRequest")
			}

		} else {
			return nil, fmt.Errorf("Missing SOCKS methods length in MethodRequest")
		}

	} else {
		return nil, fmt.Errorf("Missing SOCKS version in MethodRequest")
	}

}

func (request *MethodRequest) Bytes() []byte {

	bytes := make([]byte, 0, 2 + len(request.Methods))
	bytes = append(bytes, request.Version)
	bytes = append(bytes, byte(len(request.Methods)))
	bytes = append(bytes, request.Methods...)

	return bytes

}
