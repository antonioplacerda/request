package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// RequestAuditor audits the http requests and corresponding response.
type RequestAuditor interface {
	// SetReq sets the HTTP request to the provider to be audited.
	SetReq(req *http.Request)
	// SetResp sets the HTTP response from the provider to be audited.
	SetResp(resp *http.Response)
	// Audit sends the audit trail to the Auditor.
	Audit()
}

// Authorization adds the authorization to the http.Request.
// It should change the http.Request to include whatever needs to be added.
type Authorization interface {
	// Authorize is called after everything is set to the http.Request
	// and before making the request, so everything needed to authorize
	// it is added.
	Authorize(req *http.Request) (*http.Request, error)
}

// Request makes a new http Request with the params, headers and
// body set.
type Request struct {
	client      *http.Client
	baseURL     *url.URL
	auth        Authorization
	headers     map[string]string
	queryParams url.Values
	body        interface{}
	bodyBytes   []byte
	result      interface{}
	auditor     RequestAuditor
}

// NewRequest returns a new request.
func NewRequest(
	baseURL *url.URL,
	client *http.Client,
	accepts string,
) *Request {
	return &Request{
		baseURL: baseURL,
		client:  client,
		headers: map[string]string{
			"accept": accepts,
		},
		queryParams: make(url.Values),
	}
}

func (r *Request) WithAuthorization(a Authorization) *Request {
	r.auth = a
	return r
}

func (r *Request) WithAuditor(a RequestAuditor) *Request {
	r.auditor = a
	return r
}

func (r *Request) WithHeader(key, value string) *Request {
	if value != "" {
		r.headers[key] = value
	}
	return r
}

func (r *Request) WithHeaders(headers url.Values) *Request {
	for k, v := range headers {
		r.headers[k] = v[0]
	}
	return r
}

func (r *Request) WithResult(result interface{}) *Request {
	r.result = result
	return r
}

func (r *Request) WithJSONBody(obj interface{}) *Request {
	r.headers["Content-Type"] = "application/json"
	r.body = obj
	return r
}

func (r *Request) WithFormBody(vals url.Values) *Request {
	r.headers["Content-Type"] = "application/x-www-form-urlencoded"
	r.bodyBytes = []byte(vals.Encode())
	return r
}

// WithQParam adds a query parameter to the requests.
func (r *Request) WithQParam(key, value string) *Request {
	r.queryParams.Add(key, value)
	return r
}

func (r *Request) WithQParamPtr(key string, value *string) *Request {
	if value != nil {
		r.queryParams.Add(key, *value)
	}
	return r
}

func (r *Request) Get(ctx context.Context, path string) (interface{}, error) {
	return r.do(ctx, http.MethodGet, path)
}

func (r *Request) Post(ctx context.Context, path string) (interface{}, error) {
	if r.body != nil {
		var err error
		r.bodyBytes, err = json.Marshal(r.body)
		if err != nil {
			return nil, err
		}
	}

	return r.do(ctx, http.MethodPost, path)
}

func (r *Request) Delete(ctx context.Context, path string) (interface{}, error) {
	return r.do(ctx, http.MethodDelete, path)
}

func (r *Request) Request(ctx context.Context, method, path string) (*http.Request, error) {
	ref, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	refURL := r.baseURL.ResolveReference(ref)

	q, err := url.QueryUnescape(r.queryParams.Encode())
	if err != nil {
		return nil, err
	}
	refURL.RawQuery = q

	req, err := http.NewRequestWithContext(ctx, method, refURL.String(), bytes.NewReader(r.bodyBytes))
	if err != nil {
		return nil, err
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	if len(r.bodyBytes) > 0 {
		req.Header.Add("Content-Length", strconv.Itoa(len(r.bodyBytes)))
	}

	if r.auth != nil {
		req, err = r.auth.Authorize(req)
		if err != nil {
			return nil, err
		}
	}

	return req, err
}

func (r *Request) do(ctx context.Context, method, path string) (interface{}, error) {
	req, err := r.Request(ctx, method, path)
	if err != nil {
		return nil, err
	}

	if r.auditor != nil {
		r.auditor.SetReq(req)
		defer r.auditor.Audit()
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	if r.auditor != nil {
		r.auditor.SetResp(resp)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		re := &Error{
			HttpStatusCode: resp.StatusCode,
		}

		contentType := resp.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/html") {
			return nil, re
		}

		b, _ := io.ReadAll(resp.Body)
		re.Data = string(b)
		if contentType == "application/json" {
			json.Unmarshal(b, &re)
		}
		return nil, re
	}

	return r.decodeResult(json.NewDecoder(resp.Body))
}

func (r *Request) decodeResult(dec *json.Decoder) (interface{}, error) {
	var err error

	switch reflect.ValueOf(r.result).Kind() {
	case reflect.Ptr:
		err = dec.Decode(r.result)
	case reflect.Struct:
		t := reflect.TypeOf(r.result)
		v := reflect.New(t)
		err = dec.Decode(v.Interface())
		r.result = v.Elem().Interface()
	default:
		var v map[string]interface{}
		err = dec.Decode(&v)
		r.result = v
	}

	return r.result, err
}
