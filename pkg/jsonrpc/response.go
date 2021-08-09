package jsonrpc

import "errors"

var (
	errEmpty     = errors.New("response must have result or error")
	errAmbiguous = errors.New("error response must not have result")
)

// Result is used in successful responses.
type Result interface{}

type Response struct {
	// MUST be exactly "2.0".
	Version string `json:"jsonrpc"`

	// This member is REQUIRED on success.
	Result Result `json:"result,omitempty"`

	// This member is REQUIRED on error.
	Error *Error `json:"error,omitempty"`

	// ID is an identifier established by the Client that MUST contain a String, Number,
	// or NULL value if included.
	ID ID `json:"id,omitempty"`
}

func (r *Response) validate() error {
	if r.Result == nil && r.Error == nil {
		return errEmpty
	}

	if r.Result != nil && r.Error != nil {
		return errAmbiguous
	}

	return nil
}
