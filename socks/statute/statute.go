package statute

import (
	"errors"
)

// reply status
const (
	RepSuccess uint8 = iota
	RepServerFailure
	RepRuleFailure
	RepNetworkUnreachable
	RepHostUnreachable
	RepConnectionRefused
	RepTTLExpired
	RepCommandNotSupported
	RepAddrTypeNotSupported
	// 0x09 - 0xff unassigned
)

// error defined
var (
	ErrUnrecognizedAddrType = errors.New("unrecognized address type")
	ErrNotSupportVersion    = errors.New("not support version")
	ErrNotSupportMethod     = errors.New("not support method")
)
