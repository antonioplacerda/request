package requests

import (
	"net/http"
)

// BearerAuthorization authorizes http requests with a bearer token.
type BearerAuthorization struct {
	accessToken string
}

// NewBearerAuth returns a new BearerAuthorization.
func NewBearerAuth(accessToken string) *BearerAuthorization {
	return &BearerAuthorization{
		accessToken: accessToken,
	}
}

// Authorize sets an Authorization header with the bearer token.
func (a *BearerAuthorization) Authorize(req *http.Request) (*http.Request, error) {
	req.Header.Set("Authorization", "Bearer "+a.accessToken)
	return req, nil
}
