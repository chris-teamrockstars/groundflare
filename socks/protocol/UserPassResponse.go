package protocol

import "io"

//	+-----+--------+
//	| VER | status |
//	+-----+--------+
//	|  1  |    1   |
//	+-----+--------+

type UserPassResponse struct {
	Version byte
	Status  byte
}

func ParseUserPassResponse(reader io.Reader) (*UserPassResponse, error) {

	bytes := []byte{0, 0}

	_, err0 := io.ReadFull(reader, bytes)

	if err0 == nil {

		return &UserPassResponse{
			Version: bytes[0],
			Status:  bytes[1],
		}, nil

	} else {
		return nil, err0
	}

}
