package errors

type InvalidTopParameterError struct {
}

func (e InvalidTopParameterError) Error() string {
	return "invalid top parameter"
}
