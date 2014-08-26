package web_models

type FormError struct {
	Key   string
	Error string
}

func NewFormError(key string, err string) *FormError {
	return &FormError{Key: key, Error: err}
}
