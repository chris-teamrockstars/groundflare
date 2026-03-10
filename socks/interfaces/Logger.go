package interfaces

// Logger is used to provide debug logger
type Logger interface {
	Errorf(format string, arg ...interface{})
}
