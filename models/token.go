package models

import (
	"fmt"

	"github.com/LUSHDigital/microservice-transport-golang/config"
)

// Token - An authentication token.
type Token struct {
	Type  string `json:"type"`  // The type of auth token (e.g. JWT).
	Value string `json:"value"` // The actual token value.
}

// PrepareForRequest - Prepare a token for use with a http request.
func (t *Token) PrepareForHttp() string {
	return fmt.Sprintf("%s%s", config.AuthHeaderPrefix, t.Value)
}
