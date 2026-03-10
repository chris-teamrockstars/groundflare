package types

type Credentials map[string]string

func (credentials Credentials) Valid(user string, password string, _ string) bool {

	other, ok := credentials[user]

	if ok == true {

		if password == other {
			return true
		}

	}

	return false

}
