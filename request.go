package requests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Request request
type Request struct {
	Method  string
	URL     string
	Params  map[string]interface{}
	Headers map[string]string
	Cookies map[string]string
	Body    []byte
	Form    url.Values
	Retry   int
}

// NewRequest new request
func NewRequest(method, urlstr string, body io.Reader) *Request {
	req := &Request{
		Method:  method,
		URL:     urlstr,
		Params:  make(map[string]interface{}),
		Headers: make(map[string]string),
		Cookies: make(map[string]string),
		Form:    make(url.Values),
	}
	req.SetBody(body)
	return req
}

// SetMethod set method
func (req *Request) SetMethod(method string) *Request {
	req.Method = method
	return req
}

// SetURL set url
func (req *Request) SetURL(url string) *Request {
	req.URL = url
	return req
}

// SetParams add query args
func (req *Request) SetParams(query map[string]interface{}) *Request {
	for k, v := range query {
		req.Params[k] = v
	}
	return req
}

// SetParam params
func (req *Request) SetParam(k string, v interface{}) *Request {
	req.Params[k] = v
	return req
}

// SetBody request body
func (req *Request) SetBody(body interface{}) error {
	if body == nil {
		return nil
	}
	switch v := body.(type) {
	case string:
		req.Body = []byte(v)
	case []byte:
		req.Body = v
	case io.Reader:
		var b bytes.Buffer
		if _, err := b.ReadFrom(v); err != nil {
			return err
		}
		req.Body = b.Bytes()
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		req.Body = b
	}
	return nil
}

// SetForm set form, content-type is
func (req *Request) SetForm(k, v string) *Request {
	req.SetHeader("content-type", "application/x-www-form-urlencoded")
	req.Form.Add(k, v)
	return req
}

// SetHeader header
func (req *Request) SetHeader(k, v string) *Request {
	req.Headers[k] = v
	return req
}

// SetHeaders headers
func (req *Request) SetHeaders(kv map[string]string) *Request {
	for k, v := range kv {
		req.Headers[k] = v
	}
	return req
}

// SetCookie cookie
func (req *Request) SetCookie(k, v string) *Request {
	req.Cookies[k] = v
	return req
}

// SetCookies cookie
func (req *Request) SetCookies(kv map[string]string) *Request {
	for k, v := range kv {
		req.Cookies[k] = v
	}
	return req
}

// SetAuth base auth
func (req *Request) SetAuth(user, pass string) *Request {
	req.SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+pass)))
	return req
}

// SetRetry set retry
func (req *Request) SetRetry(retry int) *Request {
	req.Retry = retry
	return req
}

// MergeIn merge r into req
func (req *Request) MergeIn(r *Request) {
	for k, v := range r.Params {
		if _, ok := req.Params[k]; !ok {
			req.Params[k] = v
		}
	}
	for k, v := range r.Headers {
		if _, ok := req.Headers[k]; !ok {
			req.Headers[k] = v
		}
	}
	for k, v := range r.Cookies {
		if _, ok := req.Cookies[k]; !ok {
			req.Cookies[k] = v
		}
	}
	if req.Retry == 0 {
		req.Retry = r.Retry
	}
}

// Request request
func (req *Request) Request() (*http.Request, error) {
	var body io.Reader
	if len(req.Form) != 0 && len(req.Body) == 0 {
		body = strings.NewReader(req.Form.Encode())
	} else {
		body = bytes.NewReader(req.Body)
	}
	request, err := http.NewRequest(req.Method, req.URL, body)
	if err != nil {
		return nil, err
	}
	for k, v := range req.Params {
		if request.URL.RawQuery != "" {
			request.URL.RawQuery += "&"
		}
		request.URL.RawQuery += k + "=" + url.QueryEscape(fmt.Sprintf("%v", v))
	}
	for k, v := range req.Headers {
		request.Header.Set(k, v)
	}
	for k, v := range req.Cookies {
		request.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	return request, nil
}
