package microservicetransport

import (
	"io"
	"net/url"
)

// Request - Models a request to a service.
type Request struct {
	Body     io.Reader
	Method   string
	Query    url.Values
	Resource string
}
