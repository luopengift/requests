package requests

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"

	"github.com/golang-io/pkg/encoding/json"
)

// Response wrap std response
type Response struct {
	*http.Response
	once *sync.Once
	body *bytes.Buffer
}

// newResponse newResponse
func newResponse() *Response {
	return &Response{
		once: new(sync.Once),
		body: new(bytes.Buffer),
	}
}

func (resp *Response) getBody() error {
	var err error
	resp.once.Do(func() {
		if resp.Response == nil {
			return
		}
		defer resp.Response.Body.Close()
		_, err = resp.body.ReadFrom(resp.Response.Body)
	})
	return err
}

// WarpResponse warp response
func WarpResponse(resp *http.Response) *Response {
	resp2 := newResponse()
	resp2.Response = resp
	return resp2
}

// StdLib return net/http.Response
func (resp *Response) StdLib() *http.Response {
	return resp.Response
}

// Text parse parse to string
func (resp *Response) Text() (string, error) {
	if err := resp.getBody(); err != nil {
		return "", err
	}
	return resp.body.String(), nil
}

// Download parse response to a file
func (resp *Response) Download(name string) (int, error) {
	if err := resp.getBody(); err != nil {
		return 0, err
	}
	f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(resp.body.Bytes())
}

// JSON parse response
func (resp *Response) JSON(v interface{}) error {
	if err := resp.getBody(); err != nil {
		return err
	}
	return json.Unmarshal(resp.body.Bytes(), v)
}

// Dump returns the given request in its HTTP/1.x wire representation.
func (resp *Response) Dump() ([]byte, error) {
	return httputil.DumpResponse(resp.Response, true)
}

// DumpIndent warp Dump
func (resp *Response) DumpIndent() string {
	dump, _ := resp.Dump()
	var b bytes.Buffer
	for _, line := range bytes.Split(dump, []byte("\n")) {
		b.Write([]byte("< "))
		b.Write(line)
		b.WriteString("\n")
	}
	return b.String()
}

// Copy deep copy response
func (resp *Response) Copy() *Response {
	resp2 := &Response{}
	return resp2
}
