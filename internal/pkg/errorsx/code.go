package errorsx

import "net/http"

var (
	OK = &ErrorX{Code: http.StatusOK, Message: ""}

	ErrInternal = &ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError", Message: "Internal server error."}

	ErrNotFound = &ErrorX{Code: http.StatusNotFound, Reason: "NotFound", Message: "Resource not found."}
)