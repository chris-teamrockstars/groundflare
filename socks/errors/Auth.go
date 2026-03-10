package errors

import "fmt"

var (
	AuthFailure            = fmt.Errorf("Authentication failed")
	AuthUnsupportedMethod  = fmt.Errorf("Unsupported authentication method")
	AuthUnsupportedVersion = fmt.Errorf("Unsupported authentication version")
)
