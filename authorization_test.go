package requests_test

import (
	"net/http"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/antonioplacerda/requests"
)

func TestBearerAuthorization(t *testing.T) {
	c := qt.New(t)
	auth := requests.NewBearerAuth("jeff-token")

	req, err := http.NewRequest(http.MethodGet, "https://localhost", nil)
	c.Assert(err, qt.IsNil)
	c.Assert(req.Header.Get("Authorization"), qt.Equals, "")

	req, err = auth.Authorize(req)
	c.Assert(err, qt.IsNil)
	c.Assert(req.Header.Get("Authorization"), qt.Equals, "Bearer jeff-token")
}
