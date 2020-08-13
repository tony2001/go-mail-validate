package main

import (
	"github.com/tony2001/go-mail-validate/validate"
)

type Response struct {
	Valid            bool
	ValidatorResults []validate.Result
}

type ResponseError struct {
	Error string
}

func NewResponse(valid bool, results []validate.Result) *Response {
	return &Response{
		Valid:            valid,
		ValidatorResults: results,
	}
}

func NewResponseError(err string) *ResponseError {
	return &ResponseError{
		Error: err,
	}
}
