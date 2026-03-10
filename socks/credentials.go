package socks

// StaticCredentials enables using a map directly as a credential store
type StaticCredentials map[string]string

// Valid implement interface CredentialStore
func (s StaticCredentials) Valid(user, password, _ string) bool {
	pass, ok := s[user]
	return ok && password == pass
}
