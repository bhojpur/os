package httperror

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
)

var (
	Unauthorized     = ErrorCode{"Unauthorized", 401}
	PermissionDenied = ErrorCode{"PermissionDenied", 403}
	NotFound         = ErrorCode{"NotFound", 404}
	MethodNotAllowed = ErrorCode{"MethodNotAllow", 405}
	Conflict         = ErrorCode{"Conflict", 409}

	InvalidDateFormat  = ErrorCode{"InvalidDateFormat", 422}
	InvalidFormat      = ErrorCode{"InvalidFormat", 422}
	InvalidReference   = ErrorCode{"InvalidReference", 422}
	NotNullable        = ErrorCode{"NotNullable", 422}
	NotUnique          = ErrorCode{"NotUnique", 422}
	MinLimitExceeded   = ErrorCode{"MinLimitExceeded", 422}
	MaxLimitExceeded   = ErrorCode{"MaxLimitExceeded", 422}
	MinLengthExceeded  = ErrorCode{"MinLengthExceeded", 422}
	MaxLengthExceeded  = ErrorCode{"MaxLengthExceeded", 422}
	InvalidOption      = ErrorCode{"InvalidOption", 422}
	InvalidCharacters  = ErrorCode{"InvalidCharacters", 422}
	MissingRequired    = ErrorCode{"MissingRequired", 422}
	InvalidCSRFToken   = ErrorCode{"InvalidCSRFToken", 422}
	InvalidAction      = ErrorCode{"InvalidAction", 422}
	InvalidBodyContent = ErrorCode{"InvalidBodyContent", 422}
	InvalidType        = ErrorCode{"InvalidType", 422}
	ActionNotAvailable = ErrorCode{"ActionNotAvailable", 404}
	InvalidState       = ErrorCode{"InvalidState", 422}

	ServerError        = ErrorCode{"ServerError", 500}
	ClusterUnavailable = ErrorCode{"ClusterUnavailable", 503}
)

type ErrorCode struct {
	Code   string
	Status int
}

func (e ErrorCode) String() string {
	return fmt.Sprintf("%s %d", e.Code, e.Status)
}

type APIError struct {
	Code      ErrorCode
	Message   string
	Cause     error
	FieldName string
}

func NewAPIErrorLong(status int, code, message string) error {
	return NewAPIError(ErrorCode{
		Code:   code,
		Status: status,
	}, message)
}

func NewAPIError(code ErrorCode, message string) error {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

func NewFieldAPIError(code ErrorCode, fieldName, message string) error {
	return &APIError{
		Code:      code,
		Message:   message,
		FieldName: fieldName,
	}
}

// WrapFieldAPIError will cause the API framework to log the underlying err before returning the APIError as a response.
// err WILL NOT be in the API response
func WrapFieldAPIError(err error, code ErrorCode, fieldName, message string) error {
	return &APIError{
		Cause:     err,
		Code:      code,
		Message:   message,
		FieldName: fieldName,
	}
}

// WrapAPIError will cause the API framework to log the underlying err before returning the APIError as a response.
// err WILL NOT be in the API response
func WrapAPIError(err error, code ErrorCode, message string) error {
	return &APIError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

func (a *APIError) Error() string {
	if a.FieldName != "" {
		return fmt.Sprintf("%s=%s: %s", a.FieldName, a.Code, a.Message)
	}
	return fmt.Sprintf("%s: %s", a.Code, a.Message)
}

func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

func IsConflict(err error) bool {
	if apiError, ok := err.(*APIError); ok {
		return apiError.Code.Status == 409
	}

	return false
}
