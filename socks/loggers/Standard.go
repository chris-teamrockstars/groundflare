package loggers

import "log"

type Standard struct {
	*log.Logger
}

func NewStandard(logger *log.Logger) *Standard {
	return &Standard{logger}
}

// Errorf implements interface Logger
func (sf Standard) Errorf(format string, args ...interface{}) {
	sf.Printf("[E]: "+format, args...)
}
