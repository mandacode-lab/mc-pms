package common

import (
	merrmid "github.com/mandacode-com/merr/middleware"
)

// ErrorResponse represents the standard merr library error response
// This matches the exact structure returned by merr middleware
type ErrorResponse = merrmid.ErrorResponse
