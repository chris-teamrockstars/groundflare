package protocol

import "encoding/binary"
import "fmt"
import "io"
import "net"
import "groundflare/socks/types"

// +-----+-----+-------+------+----------+-----------+
// | VER | REP |  RSV  | ATYP | BND.ADDR | BND].PORT |
// +-----+-----+-------+------+----------+-----------+
// |  1  |  1  | X'00' |  1   | Variable |    2      |
// +-----+-----+-------+------+----------+-----------+
type ConnectionResponse struct {
	Version byte
	Reply   byte
	Address types.Address
}

func ParseConnectionResponse(reader io.Reader) (*ConnectionResponse, error) {

	version := []byte{0x00}
	_, err0 := reader.Read(version)

	if err0 == nil && version[0] == Version5 {

		reply := []byte{0x00}
		_, err1 := reader.Read(reply)

		if err1 == nil {

			reserved := []byte{0xff}
			_, err2  := reader.Read(reserved)

			if err2 == nil && reserved[0] == 0x00 {

				address_type := []byte{0x00}
				_, err3      := reader.Read(address_type)

				if err3 == nil {

					switch types.AddressType(address_type[0]) {
					case types.AddressTypeIPv4:

						address := make([]byte, net.IPv4len + 2)
						_, err4 := io.ReadAtLeast(reader, address, net.IPv4len + 2)

						if err4 == nil {

							ip   := net.IPv4(address[0], address[1], address[2], address[3])
							port := int(binary.BigEndian.Uint16(address[net.IPv4len:]))

							return &ConnectionResponse{
								Version: version[0],
								Reply:   reply[0],
								Address: types.Address{
									Type: types.AddressTypeIPv4,
									IP:   ip, // Go uses internal 16 bytes IPv6 format
									Port: port,
								},
							}, nil

						} else {

							hex_str := ""

							for _, byt := range address {
								hex_str += fmt.Sprintf("0x%02x ", byt)
							}

							return nil, fmt.Errorf("Unsupported address %s in ConnectionResponse", hex_str)

						}

					case types.AddressTypeDomain:

						address_length := []byte{0x00}
						_, err4        := reader.Read(address_length)

						if err4 == nil && int(address_length[0]) > 0 {

							address := make([]byte, int(address_length[0]) + 2)
							_, err5 := io.ReadAtLeast(reader, address, int(address_length[0]) + 2)

							if err5 == nil {

								fqdn := string(address[0:int(address_length[0])])
								port := int(binary.BigEndian.Uint16(address[int(address_length[0]):]))

								return &ConnectionResponse{
									Version: version[0],
									Reply:   reply[0],
									Address: types.Address{
										Type: types.AddressTypeDomain,
										FQDN: fqdn,
										Port: port,
									},
								}, nil

							} else {

								hex_str := ""
								for _, byt := range address {
									hex_str += fmt.Sprintf("0x%02x ", byt)
								}

								return nil, fmt.Errorf("Unsupported address %s in ConnectionResponse", hex_str)

							}

						} else {
							return nil, fmt.Errorf("Unsupported address length in ConnectionResponse")
						}

					case types.AddressTypeIPv6:

						address := make([]byte, net.IPv6len + 2)
						_, err4 := io.ReadAtLeast(reader, address, net.IPv6len + 2)

						if err4 == nil {

							ip   := net.IP(address[0:net.IPv6len])
							port := int(binary.BigEndian.Uint16(address[net.IPv6len:]))

							return &ConnectionResponse{
								Version: version[0],
								Reply:   reply[0],
								Address: types.Address{
									Type: types.AddressTypeIPv6,
									IP:   ip,
									Port: port,
								},
							}, nil

						} else {

							hex_str := ""

							for _, byt := range address {
								hex_str += fmt.Sprintf("0x%02x ", byt)
							}

							return nil, fmt.Errorf("Unsupported address %s in ConnectionResponse", hex_str)

						}

					default:
						return nil, fmt.Errorf("Unsupported SOCKS address type 0x%02x in ConnectionResponse", address_type[0])
					}

				} else {
					return nil, fmt.Errorf("Missing SOCKS address type in ConnectionResponse")
				}

			} else {
				return nil, fmt.Errorf("Missing SOCKS reserved byte in ConnectionResponse")
			}

		} else {
			return nil, fmt.Errorf("Unsupported SOCKS reply 0x%02x in ConnectionResponse", reply[0])
		}

	} else {
		return nil, fmt.Errorf("Unsupported SOCKS version \"%d\" in ConnectionResponse", int(version[0]))
	}

}

func (response *ConnectionResponse) Bytes() []byte {

	bytes := make([]byte, 0)
	bytes = append(bytes, response.Version)
	bytes = append(bytes, response.Reply)
	bytes = append(bytes, byte(0x00)) // RSV
	bytes = append(bytes, response.Address.Bytes()...)

	return bytes

}

