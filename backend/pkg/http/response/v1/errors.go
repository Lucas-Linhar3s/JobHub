package v1

import (
	"fmt"
)

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// more biz errors
	ErrEmailAlreadyUse = newError(1001, "The email is already in use.")
	ErrEmailNotExists  = newError(1002, "The email does not exist.")
	ErrInvalidPassword = newError(1003, "Invalid password.")
	ErrInvalidState    = newError(1004, "Invalid state.")
	ErrStructureError  = newError(1002, "Error in the structure.")
)

// CheckError checks error
func CheckError(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %s", msg, err)
	}
	return nil
}
