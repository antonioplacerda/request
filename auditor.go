package requests

import (
	"net/http"
	"net/http/httputil"
)

type Auditor interface {
	// UserSuccess audits the requests and response when the http status
	// code is <300.
	UserSuccess(userID, action string, metadata map[string]interface{})
	// UserFail audits the requests and response when the http status
	// code is >=300.
	UserFail(userID, action string, metadata map[string]interface{})
}

type reqAuditor struct {
	auditor Auditor
	req     []byte
	resp    []byte
	success bool
	action  string
	userID  string
}

func NewReqAuditor(auditor Auditor, userID, action string) *reqAuditor {
	return &reqAuditor{
		auditor: auditor,
		action:  action,
		userID:  userID,
	}
}

// SetReq sets the HTTP request to the provider to be audited.
func (a *reqAuditor) SetReq(req *http.Request) {
	d, _ := httputil.DumpRequestOut(req, true)
	a.req = d
}

// SetResp sets the HTTP response from the provider to be audited.
func (a *reqAuditor) SetResp(resp *http.Response) {
	a.success = resp.StatusCode < 300

	d, _ := httputil.DumpResponse(resp, true)
	a.resp = d
}

// Audit sends the audit trail to the Auditor.
func (a *reqAuditor) Audit() {
	if a == nil {
		return
	}
	md := map[string]interface{}{
		"req":  a.req,
		"resp": a.resp,
	}
	if a.success {
		a.auditor.UserSuccess(a.userID, a.action, md)
		return
	}
	a.auditor.UserFail(a.userID, a.action, md)
}
