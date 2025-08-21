package apierror

type APIerror struct {
	code int
	err  error
}

func (e *APIerror) Error() string {
	return e.err.Error()
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
