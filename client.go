package microservicetransport

import (
	"net/http"
	"time"
)

// DefaultHttpClient - returns a default http.Client implementation
func DefaultHttpClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}
