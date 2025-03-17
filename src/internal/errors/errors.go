package errors

import (
	"github.com/go-kratos/kratos/v2/errors"
)

// Error codes
var (
	ErrorInvalidName = errors.New(400, "INVALID_NAME", "name is invalid or empty")
	ErrorNameTooLong = errors.New(400, "NAME_TOO_LONG", "name is too long (max 50 characters)")
	ErrorServerError = errors.New(500, "SERVER_ERROR", "internal server error")
)
