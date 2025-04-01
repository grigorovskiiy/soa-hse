package errors

type PostNotFoundError struct {
}

func (err PostNotFoundError) Error() string {
	return "Post not found"
}
