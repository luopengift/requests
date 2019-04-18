package requests

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptrace"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// var
var (
	ErrEmptyProxy = errors.New("proxy is empty")
)

// Session httpclient session
// Clients and Transports are safe for concurrent use by multiple goroutines
// for efficiency should only be created once and re-used.
// so, session is also safe for concurrent use by multiple goroutines.
type Session struct {
	*http.Transport
	*http.Client
	option  *Request
	LogFunc func(string, ...interface{})
	errs    chan error
}

// New new session
func New() *Session {
	tr := &http.Transport{
		MaxIdleConns:        10,
		TLSHandshakeTimeout: 10 * time.Second,
		IdleConnTimeout:     120 * time.Second,
		DisableCompression:  true,
		DisableKeepAlives:   false,
	}
	jar, _ := cookiejar.New(nil)
	return &Session{
		Transport: tr,
		Client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: tr,
			Jar:       jar,
		},
		option: NewRequest("", "", nil),
		LogFunc: func(format string, v ...interface{}) {
			fmt.Fprintf(os.Stderr, format+"\n", v...)
		},
		errs: make(chan error),
	}
}

// SetParam query
func (sess *Session) SetParam(k string, v interface{}) *Session {
	sess.option.SetParam(k, v)
	return sess
}

// SetParams querys
func (sess *Session) SetParams(kv map[string]interface{}) *Session {
	for k, v := range kv {
		sess.option.SetParam(k, v)
	}
	return sess
}

// SetHeader header
func (sess *Session) SetHeader(k, v string) *Session {
	sess.option.SetHeader(k, v)
	return sess
}

// SetHeaders headers
func (sess *Session) SetHeaders(kv map[string]string) *Session {
	for k, v := range kv {
		sess.option.SetHeader(k, v)
	}
	return sess
}

// SetCookie cookie
func (sess *Session) SetCookie(k, v string) *Session {
	sess.option.SetCookie(k, v)
	return sess
}

// SetAuth base auth, proxy is Proxy-Authorization
func (sess *Session) SetAuth(user, pass string) *Session {
	sess.option.SetAuth(user, pass)
	return sess
}

// SetProxy set proxy addr
// os.Setenv("HTTP_PROXY", "http://127.0.0.1:9743")
// os.Setenv("HTTPS_PROXY", "https://127.0.0.1:9743")
func (sess *Session) SetProxy(addr string) error {
	if addr == "" {
		return ErrEmptyProxy
	}
	proxyURL, err := url.Parse(addr)
	if err != nil {
		return err
	}
	switch proxyURL.Scheme {
	case "socks5", "socks4":
		dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, nil, proxy.Direct)
		if err != nil {
			return err
		}
		sess.Transport.Dial = dialer.Dial
	default:
		sess.Transport.Proxy = http.ProxyURL(proxyURL)
	}
	return nil
}

// SetLogFunc set log handler
func (sess *Session) SetLogFunc(f func(string, ...interface{})) *Session {
	sess.LogFunc = f
	return sess
}

// SetRetry set retry times if request fail
func (sess *Session) SetRetry(retry int, retryFunc ...func(Response, error) error) *Session {
	sess.option.SetRetry(retry)
	return sess
}

// SetTimeout set client timeout
func (sess *Session) SetTimeout(timeout int) *Session {
	sess.Client.Timeout = time.Duration(timeout) * time.Second
	return sess
}

// SetKeepAlives set transport disableKeepAlives default transport is keepalive,
// if set false, only use the connection to the server for a single HTTP request.
func (sess *Session) SetKeepAlives(keepAlives bool) *Session {
	sess.Transport.DisableKeepAlives = !keepAlives
	return sess
}

// DoRequest send a request and return a response
func (sess *Session) DoRequest(request *Request, ctx ...context.Context) (*Response, error) {
	request.MergeIn(sess.option)
	var req *http.Request
	var resp *http.Response
	var err error
	for i := 0; i <= request.Retry; i++ {
		req, err = request.Request()
		if err != nil {
			return nil, err
		}
		if len(ctx) != 0 {
			req = req.WithContext(ctx[0]) // !!! WithContext returns a shallow copy of r with its context changed to ctx
		}
		if resp, err = sess.Client.Do(req); err == nil {
			break
		}
		if i != 0 {
			sess.LogFunc("retry[%d/%d], err=%v", i, request.Retry, err)
		}
	}
	return WarpResponse(resp), err
}

// Do http request
func (sess *Session) Do(method, url, contentType string, body io.Reader) (*Response, error) {
	req := NewRequest(method, url, body).SetHeader("Content-Type", contentType)
	return sess.DoRequest(req)
}

// DoWithContext http request
func (sess *Session) DoWithContext(ctx context.Context, method, url, contentType string, body io.Reader) (*Response, error) {
	req := NewRequest(method, url, body).SetHeader("Content-Type", contentType)
	return sess.DoRequest(req, ctx)
}

// Get send get request
func (sess *Session) Get(url string) (*Response, error) {
	return sess.Do("GET", url, "", nil)
}

// GetWithContext http request
func (sess *Session) GetWithContext(ctx context.Context, url string) (*Response, error) {
	return sess.DoWithContext(ctx, "GET", url, "", nil)
}

// Post send post request
func (sess *Session) Post(url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.Do("POST", url, contentType, body)
}

// PostWithContext send post request
func (sess *Session) PostWithContext(ctx context.Context, url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.DoWithContext(ctx, "POST", url, contentType, body)
}

// PostForm post form request
func (sess *Session) PostForm(url string, data url.Values) (resp *Response, err error) {
	return sess.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// PostFormWithContext post form request
func (sess *Session) PostFormWithContext(ctx context.Context, url string, data url.Values) (resp *Response, err error) {
	return sess.PostWithContext(ctx, url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// Put send put request
func (sess *Session) Put(url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.Do("PUT", url, contentType, body)
}

// PutWithContext send put request
func (sess *Session) PutWithContext(ctx context.Context, url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.DoWithContext(ctx, "PUT", url, contentType, body)
}

// Delete send delete request
func (sess *Session) Delete(url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.Do("DELETE", url, contentType, body)
}

// DeleteWithContext send delete request
func (sess *Session) DeleteWithContext(ctx context.Context, url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.DoWithContext(ctx, "DELETE", url, contentType, body)
}

// DebugTrace trace a request
// 	sess := requests.New()
//	req, err := requests.NewRequest("GET", "http://www.baidu.com", nil)
//	if err != nil {
//		sess.LogFunc("%v", err)
//	}
//	sess.DebugTrace(req)
//
// * Connect: www.baidu.com:80
// * Resolved Host: www.baidu.com
// * Resolved DNS: [14.215.177.39 14.215.177.38], Coalesced: false, err=<nil>
// *   Trying tcp 14.215.177.39:80...
// * Completed connection: tcp 14.215.177.39:80, err=<nil>
// * Got Conn: {0xc0000ac020 false false 0s}
// > GET / HTTP/1.1
// > Host: www.baidu.com
// > User-Agent: Go-http-client/1.1
// > Accept-Encoding: gzip
// >
// >
// < HTTP/1.1 200 OK
// < Transfer-Encoding: chunked
// < Bdpagetype: 1
// < Bdqid: 0x85d6ecb5000fcc70
// ...more...
func (sess *Session) DebugTrace(request *Request) {
	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			sess.LogFunc("* Connect: %v", hostPort)
		},
		ConnectStart: func(network, addr string) {
			sess.LogFunc("*   Trying %v %v...", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			sess.LogFunc("* Completed connection: %v %v, err=%v", network, addr, err)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			sess.LogFunc("* Got Conn: %v -> %v", connInfo.Conn.LocalAddr(), connInfo.Conn.RemoteAddr())
		},
		DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
			sess.LogFunc("* Resolved Host: %v", dnsInfo.Host)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			var ipaddrs []string
			for _, ipaddr := range dnsInfo.Addrs {
				ipaddrs = append(ipaddrs, ipaddr.String())
			}
			sess.LogFunc("* Resolved DNS: %v, Coalesced: %v, err=%v", ipaddrs, dnsInfo.Coalesced, dnsInfo.Err)
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			sess.LogFunc("* SSL HandshakeComplete: %v", state.HandshakeComplete)
		},
		WroteRequest: func(reqInfo httptrace.WroteRequestInfo) {
		},
	}
	fmt.Println(trace)
	req, err := request.Request()
	if err != nil {
		sess.LogFunc("new request: %v", err)
		return
	}
	ctx := httptrace.WithClientTrace(req.Context(), trace)
	req2 := req.WithContext(ctx)
	resp, err := sess.Transport.RoundTrip(req2)
	sess.LogFunc(DumpRequestIndent(req2))
	if err != nil {
		sess.LogFunc("response error: %v", err)
		return
	}
	sess.LogFunc(WarpResponse(resp).DumpIndent())
}
