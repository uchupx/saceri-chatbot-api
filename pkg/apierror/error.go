package apierror

import "fmt"

type APIerror struct {
	code int
	err  error
}

func (e *APIerror) Error() string {
	return e.err.Error()
}

func (e *APIerror) Unwrap() error {
	return e.err
}

func (e *APIerror) Code() int {
	return e.code
}

func NewAPIError(code int, err error) *APIerror {
	return &APIerror{
		code: code,
		err:  err,
	}
}

var ERR_NOT_FOUND error = fmt.Errorf("data is not found")
