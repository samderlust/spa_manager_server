package resterrors

import "net/http"

//RestError describe the error that this service returns
type RestError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

//NewBadRequestError return a pointer to RestError with bad request type
func NewBadRequestError(message string) *RestError {
	return &RestError{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "badRequestError",
	}
}

//NewNotFoundError return a pointer to RestError with bad request type
func NewNotFoundError(message string) *RestError {
	return &RestError{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "notFoundError",
	}
}

//NewInternalServerError return a pointer to RestError with internal server error
func NewInternalServerError(message string) *RestError {
	return &RestError{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "internalServerError",
	}
}
