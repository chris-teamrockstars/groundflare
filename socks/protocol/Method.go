package protocol

const (
	MethodNoAuth              = byte(0x00)
	MethodGSSAPIAuth          = byte(0x01) // TODO: GSSAPI is unsupported for now
	MethodUserPassAuth        = byte(0x02)
	MethodNoAcceptableMethods = byte(0xff)
)
