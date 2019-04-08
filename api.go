package requests

import (
	"net/url"
)

var sess = New()

// Get send get request
func Get(url string) (*Response, error) {
	req := NewRequest("GET", url, nil)
	return sess.DoRequest(req)
}

// Post send post request
func Post(url, contentType string, body interface{}) (*Response, error) {
	req := NewRequest("POST", url, nil).SetHeader("Content-Type", contentType)
	if err := req.SetBody(body); err != nil {
		return nil, err
	}
	return sess.DoRequest(req)

}

// PUT send post request
func PUT(url string, body interface{}) (resp *Response, err error) {
	req := NewRequest("PUT", url, nil)
	if err := req.SetBody(body); err != nil {
		return nil, err
	}
	return sess.DoRequest(req)
}

// Delete send post request
func Delete(url string, body interface{}) (resp *Response, err error) {
	req := NewRequest("DELETE", url, nil)
	if err := req.SetBody(body); err != nil {
		return nil, err
	}
	return sess.DoRequest(req)

}

// Head send post request
func Head(url string) (resp *Response, err error) {
	req := NewRequest("HEAD", url, nil)
	return sess.DoRequest(req)
}

// PostForm send post request,  content-type = application/x-www-form-urlencoded
func PostForm(url string, data url.Values) (resp *Response, err error) {
	return sess.PostForm(url, data)
}

// Wget download a file from remote.
func Wget(url, name string) (int, error) {
	resp, err := sess.Get(url)
	if err != nil {
		return 0, err
	}
	return resp.Download(name)
}
