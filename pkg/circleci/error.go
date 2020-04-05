package circleci

// ErrorResponse defines struct for api response of error.
type ErrorResponse struct {
	Message string `json:"message"`
}

// String implement interface of Stringer.
func (e ErrorResponse) String() string {
	return e.Message
}

// ErrorResponse implement interface of error.
func (e ErrorResponse) Error() string {
	return e.String()
}

// IsEmpty check whether error is empty.
func (e ErrorResponse) IsEmpty() bool {
	return e.String() == ""
}
