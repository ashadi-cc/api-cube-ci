package main

import "fmt"

type errorData struct {
	ID      string      `json:"id,omitempty"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func (e *errorData) Error() string {
	return fmt.Sprintf("%s: %s %+v", e.ID, e.Message, e.Details)
}

var apiErrors = map[string]string{
	errInternalServerError: "Unexpected server error.",
	errMethodNotAllowed:    "Method not allowed",
	errBadRequest:          "Bad request",
}

func newError(ID string, errorMessages map[string]string, details ...map[string]string) *errorData {
	err := &errorData{}
	err.ID = ID
	err.Message = fmt.Sprintf(errorMessages[ID])
	if len(details) > 0 {
		err.Details = details
	}

	return err
}

func newAPIError(ID string, details ...map[string]string) *errorData {
	return newError(ID, apiErrors, details...)
}
