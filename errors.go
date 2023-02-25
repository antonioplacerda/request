package requests

import (
	"errors"
	"fmt"
)

// Error handles the errors received from the provider.
type Error struct {
	HttpStatusCode int    `json:"-"`
	Data           string `json:"-"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("error: %d, %s", e.HttpStatusCode, e.Data)
}

func getError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return nil
}
