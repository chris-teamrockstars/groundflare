package protocol

import "encoding/binary"
import "fmt"
import "io"
import "net"
import "groundflare/socks/types"

// +-----+-----+-------+------+----------+----------+
// | VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
// +-----+-----+-------+------+----------+----------+
// |  1  |  1  | X'00' |  1   | Variable |    2     |
// +-----+-----+-------+------+----------+----------+
type ConnectionRequest struct {
	Version byte          // Version4, Version5
	Command byte          // CommandConnect, CommandBind, CommandAssociate
	Address types.Address // Destination
}

func ParseConnectionRequest(reader io.Reader) (*ConnectionRequest, error) {

	version := []byte{0x00}
	_, err0 := reader.Read(version)

	if err0 == nil && version[0] == Version5 {

		command := []byte{0x00}
		_, err1 := reader.Read(command)

		if err1 == nil && (
			command[0] == CommandConnect ||
			command[0] == CommandBind ||
			command[0] == CommandAssociate) {

			reserved := []byte{0xff}
			_, err2 := reader.Read(reserved)

			if err2 == nil && reserved[0] == 0x00 {

				address_type := []byte{0x00}
				_, err3 := reader.Read(address_type)

				if err3 == nil {

					switch types.AddressType(address_type[0]) {
					case types.AddressTypeIPv4:

						address := make([]byte, net.IPv4len + 2)
						_, err4 := io.ReadAtLeast(reader, address, net.IPv4len + 2)

						if err4 == nil {

							ip   := net.IPv4(address[0], address[1], address[2], address[3])
							port := int(binary.BigEndian.Uint16(address[net.IPv4len:]))

							return &ConnectionRequest{
								Version: version[0],
								Command: command[0],
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

							return nil, fmt.Errorf("Unsupported address %s in ConnectionRequest", hex_str)

						}

					case types.AddressTypeBrokenClientAuthentication:

						return nil, fmt.Errorf("Unsupported SOCKS address type 0x%02x in ConnectionRequest", address_type[0])

					case types.AddressTypeDomain:

						address_length := []byte{0x00}
						_, err4        := reader.Read(address_length)

						if err4 == nil && int(address_length[0]) > 0 {

							address := make([]byte, int(address_length[0]) + 2)
							_, err5 := io.ReadAtLeast(reader, address, int(address_length[0]) + 2)

							if err5 == nil {

								fqdn := string(address[0:int(address_length[0])])
								port := int(binary.BigEndian.Uint16(address[int(address_length[0]):]))

								return &ConnectionRequest{
									Version: version[0],
									Command: command[0],
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

								return nil, fmt.Errorf("Unsupported address %s in ConnectionRequest", hex_str)

							}

						} else {
							return nil, fmt.Errorf("Unsupported address length in ConnectionRequest")
						}

					case types.AddressTypeIPv6:

						address := make([]byte, net.IPv6len + 2)
						_, err4 := io.ReadAtLeast(reader, address, net.IPv6len + 2)

						if err4 == nil {

							ip   := net.IP(address[0:net.IPv6len])
							port := int(binary.BigEndian.Uint16(address[net.IPv6len:]))

							return &ConnectionRequest{
								Version: version[0],
								Command: command[0],
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

							return nil, fmt.Errorf("Unsupported address %s in ConnectionRequest", hex_str)

						}

					default:
						return nil, fmt.Errorf("Unsupported SOCKS address type 0x%02x in ConnectionRequest", address_type[0])
					}

				} else {
					return nil, fmt.Errorf("Missing SOCKS address type in ConnectionRequest")
				}

			} else {
				return nil, fmt.Errorf("Missing SOCKS reserved byte in ConnectionRequest")
			}

		} else {
			return nil, fmt.Errorf("Unsupported SOCKS command 0x%02x in ConnectionRequest", command[0])
		}

	} else if err0 == nil && version[0] == Version4 {
		return nil, fmt.Errorf("Unsupported SOCKS version \"4\" in ConnectionRequest")
	} else {
		return nil, fmt.Errorf("Unsupported SOCKS version \"%d\" in ConnectionRequest", int(version[0]))
	}

}

func (request *ConnectionRequest) Bytes() []byte {

	bytes := make([]byte, 0)
	bytes = append(bytes, request.Version)
	bytes = append(bytes, request.Command)
	bytes = append(bytes, byte(0x00)) // RSV
	bytes = append(bytes, request.Address.Bytes()...)

	return bytes

}
