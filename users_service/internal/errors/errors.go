package errors

type LoginError struct {
}

func (err LoginError) Error() string {
	return "Login error"
}

type AlreadyRegisteredError struct {
}

func (err AlreadyRegisteredError) Error() string {
	return "User with this login already exists"
}
