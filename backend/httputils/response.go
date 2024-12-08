package httputils

import (
	"time"

	"github.com/johannfh/go-utils/assert"
)

type responseStatus string

const (
	// The response status in case of an error
	responseStatusError responseStatus = "error"

	// The response status in case of success
	responseStatusSuccess responseStatus = "success"
)

type Response[T any] struct {
	// may be either "error" or "success"
	Status     responseStatus `json:"status"`
	StatusCode int            `json:"statusCode"`

	// empty if [Status] is "success"
	Error *ResponseError `json:"error,omitempty"`

	// empty if [Status] is "error"
	Data *T `json:"data,omitempty"`
}

type NewResponseOpt[T any] func(*Response[T])

func NewResponse[T any](options ...NewResponseOpt[T]) *Response[T] {
	response := &Response[T]{}

	for _, option := range options {
		option(response)
	}

	assert.NotEmpty(string(response.Status), "missing response status", "response", response)

	return response
}

func WithStatusError[T any](response *Response[T]) {
	response.Status = responseStatusError
}

func WithStatusSuccess[T any](response *Response[T]) {
	response.Status = responseStatusSuccess
}

func WithStatusCode[T any](code int) NewResponseOpt[T] {
	return func(r *Response[T]) {
		r.StatusCode = code
	}
}

func WithResponseError[T any](err *ResponseError) NewResponseOpt[T] {
	return func(r *Response[T]) {
		WithStatusError(r)
		r.Error = err
	}
}

func WithResponseData[T any](data *T) NewResponseOpt[T] {
	return func(r *Response[T]) {
		WithStatusSuccess(r)
		r.Data = data
	}
}

type ResponseError struct {
	// Example: "RESOURCE_NOT_FOUND"
	Code string `json:"code"` // required

	// Example: "The requested resource was not found."
	Message string `json:"message"` // required

	// Example: "The user with the ID '12345' does not exist in our records."
	Details string `json:"details"` // optional

	// Example: "/api/v1/users/12345"
	Path string `json:"path"` // optional

	// Example: "1970-1-1T12:30:45Z"
	Timestamp string `json:"timestamp"` // optional
}

type NewResponseErrorOpt func(*ResponseError)

func NewResponseError(options ...NewResponseErrorOpt) *ResponseError {
	responseError := &ResponseError{}

	for _, option := range options {
		option(responseError)
	}

	return responseError
}

func WithErrorCode(code string) NewResponseErrorOpt {
	return func(re *ResponseError) {
		re.Code = code
	}
}

func WithErrorMessage(msg string) NewResponseErrorOpt {
	return func(re *ResponseError) {
		re.Message = msg
	}
}

func WithErrorDetails(details string) NewResponseErrorOpt {
	return func(re *ResponseError) {
		re.Details = details
	}
}

func WithErrorPath(path string) NewResponseErrorOpt {
	return func(re *ResponseError) {
		re.Path = path
	}
}

func WithErrorTimestamp(timestamp time.Time) NewResponseErrorOpt {
	return func(re *ResponseError) {
		re.Timestamp = timestamp.UTC().Format(time.RFC3339)
	}
}
