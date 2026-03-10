package statute

import "io"

// UserPassRequest is the negotiation user's password request packet
// The SOCKS handshake user's password request is formed as follows:
//

// UserPassReply is the negotiation user's password reply packet
// The SOCKS handshake user's password response is formed as follows:
//
//	+-----+--------+
//	| VER | status |
//	+-----+--------+
//	|  1  |     1  |
//	+-----+--------+
type UserPassReply struct {
	Ver    byte
	Status byte
}

// ParseUserPassReply parse user's password reply packet.
func ParseUserPassReply(r io.Reader) (upr UserPassReply, err error) {
	bb := []byte{0, 0}
	if _, err = io.ReadFull(r, bb); err != nil {
		return
	}
	upr.Ver, upr.Status = bb[0], bb[1]
	return
}
