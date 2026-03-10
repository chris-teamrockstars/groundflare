package interfaces

type Credentials interface {
	Valid(user, password, userAddr string) bool
}
