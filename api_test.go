package requests_test

import (
	"testing"

	"github.com/luopengift/requests"
)

func Test_apiGet(t *testing.T) {
	resp, err := requests.Get("http://httpbin.org/get")
	t.Log(resp.DumpIndent(), err)
}

func Test_apiPost(t *testing.T) {
	resp, err := requests.Post("http://httpbin.org/post", "application/json", `{"a":"b"}`)
	t.Log(resp.DumpIndent(), err)
}
