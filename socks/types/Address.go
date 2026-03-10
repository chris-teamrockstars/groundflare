package types

import "fmt"
import "net"
import "strconv"

type Address struct {
	Type AddressType
	FQDN string
	IP   net.IP
	Port int
}

func ParseAddress(raw string) (*Address, error) {

	host, port, err0 := net.SplitHostPort(raw)

	if err0 == nil {

		port, err1 := strconv.Atoi(port)

		if err1 == nil {

			ip := net.ParseIP(host)

			if ipv4 := ip.To4(); ipv4 != nil {

				return &Address{
					Type: AddressTypeIPv4,
					FQDN: "",
					IP:   ip,
					Port: port,
				}, nil

			} else if ipv6 := ip.To16(); ipv6 != nil {

				return &Address{
					Type: AddressTypeIPv6,
					FQDN: "",
					IP:   ip,
					Port: port,
				}, nil

			} else if len(host) > 0 {

				return &Address{
					Type: AddressTypeDomain,
					FQDN: host,
					IP:   nil,
					Port: port,
				}, nil

			} else {
				return nil, fmt.Errorf("Unsupported address format")
			}

		} else {
			return nil, err1
		}

	} else {
		return nil, err0
	}

}

func (address *Address) Bytes() []byte {

	bytes := make([]byte, 0)

	switch address.Type {
	case AddressTypeIPv4:

		bytes = append(bytes, byte(address.Type))
		bytes = append(bytes, address.IP.To4()...)
		bytes = append(bytes, byte(address.Port >> 8))
		bytes = append(bytes, byte(address.Port))

	case AddressTypeDomain:

		bytes = append(bytes, byte(address.Type))
		bytes = append(bytes, byte(len(address.FQDN)))
		bytes = append(bytes, []byte(address.FQDN)...)
		bytes = append(bytes, byte(address.Port >> 8))
		bytes = append(bytes, byte(address.Port))

	case AddressTypeIPv6:

		bytes = append(bytes, byte(address.Type))
		bytes = append(bytes, address.IP.To16()...)
		bytes = append(bytes, byte(address.Port >> 8))
		bytes = append(bytes, byte(address.Port))


	}

	return bytes

}

func (address *Address) String() string {

	if len(address.IP) > 0 {
		return net.JoinHostPort(address.IP.String(), strconv.Itoa(address.Port))
	} else if len(address.FQDN) > 0 {
		return net.JoinHostPort(address.FQDN, strconv.Itoa(address.Port))
	} else {
		return "0.0.0.0:0"
	}

}
