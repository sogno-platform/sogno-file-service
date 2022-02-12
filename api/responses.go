// SPDX-License-Identifier: Apache-2.0

package api

import (
	"time"
)

type ResponseErrorData struct {
	// HTTP response status code
	Code int `json:"code" validate:"required"`
	// Description of the error
	Message string `json:"message" validate:"required"`
}

// @Description An error
type ResponseError struct {
	Error ResponseErrorData `json:"error" validate:"required"`
}

type ResponseFileData struct {
	// ID of file
	FileID string `json:"fileID" validate:"required"`
	// Last modified timestamp of file
	LastModified time.Time `json:"lastModified" validate:"required"`
	// URL of file
	URL string `json:"url,omitempty"`
}

// @Description A single file
type ResponseFile struct {
	Data ResponseFileData `json:"data" validate:"required"`
}

// @Description Multiple files (NOTE: Partial responses are possible.
// @Description  In this case `data` and `error` will be returned.)
type ResponseFiles struct {
	Data  []ResponseFileData `json:"data,omitempty"`
	Error *ResponseErrorData `json:"error,omitempty"`
}

// @Description Empty successful response
type ResponseEmpty struct {
	Data struct{} `json:"data" validate:"required"`
}
