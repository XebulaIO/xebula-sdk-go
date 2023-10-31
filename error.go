package xebula

import "fmt"

type XebulaError struct {
	Message string `json:"message"`
}

type InvalidURLError struct {
	XebulaError
	Details string
}

type NotImplementedError struct {
	XebulaError
}

func (xe XebulaError) Error() string {
	return fmt.Sprintf("%s", xe.Message)
}

func NewError(_type string, args ...interface{}) error {
	switch _type {
	case INVALID_URL_ERROR:
		if len(args) > 0 {
			return InvalidURLError{
				XebulaError{Message: "invalid API url"},
				args[0].(string),
			}

		}
		return InvalidURLError{
			XebulaError{Message: "invalid API url"},
			"",
		}
	case NOT_IMPLEMENTED_ERROR:
		return &NotImplementedError{
			XebulaError{Message: "method not implemented"},
		}
	}
	if len(args) > 0 {
		return XebulaError{Message: args[0].(string)}
	}
	return XebulaError{}
}
