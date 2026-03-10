package types

type AddressType byte

const (
	AddressTypeIPv4   AddressType = 0x01
	AddressTypeDomain AddressType = 0x03
	AddressTypeIPv6   AddressType = 0x04

	// XXX: This happens for broken SOCKS clients
	AddressTypeBrokenClientAuthentication AddressType = 0x02
)
