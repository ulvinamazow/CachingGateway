package model

import "net/http"

// to keep the status codes, headers, body
type CachedItem struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}
