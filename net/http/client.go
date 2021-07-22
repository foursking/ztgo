package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	xhttp "net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go/ext"

	"git.code.oa.com/qdgo/core/config/env"
	"git.code.oa.com/qdgo/core/net/util/breaker"

	"github.com/gogo/protobuf/proto"
	pkgerr "github.com/pkg/errors"
)

const (
	_minRead = 64 * 1024 // 64kb
)

var (
	_userAgent = "qdgo_http "
)

func init() {
	n, err := os.Hostname()
	if err == nil {
		_userAgent = _userAgent + runtime.Version() + " " + n
	}
}

// Client is http client.
type Client struct {
	opts      *ClientOptions
	client    *xhttp.Client
	dialer    *net.Dialer
	transport xhttp.RoundTripper

	urlConf  map[string]*ClientOptions
	hostConf map[string]*ClientOptions
	mu       sync.RWMutex
	breaker  *breaker.Group
}

// NewClient new a http client.
func NewClient(opts ...ClientOption) *Client {
	client := new(Client)
	options := DefaultClientOptions
	for _, o := range opts {
		o(&options)
	}
	client.opts = &options
	client.dialer = &net.Dialer{
		Timeout:   options.DialTimeout,
		KeepAlive: options.KeepAlive,
	}
	client.transport = &xhttp.Transport{
		DialContext:     client.dialer.DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.client = &xhttp.Client{
		Transport: client.transport,
	}
	client.urlConf = make(map[string]*ClientOptions)
	client.hostConf = make(map[string]*ClientOptions)
	client.breaker = breaker.NewGroup(options.Breaker)
	if options.Timeout <= 0 {
		panic("must config http timeout!!!")
	}
	for u, cfg := range options.URL {
		client.urlConf[u] = cfg
	}
	for host, cfg := range options.Host {
		client.hostConf[host] = cfg
	}
	return client
}

// SetTransport set client transport
func (client *Client) SetTransport(t xhttp.RoundTripper) {
	client.transport = t
	client.client.Transport = t
}

// SetConfig set client config.
func (client *Client) SetConfig(c *ClientOptions) {
	client.mu.Lock()
	if c.Timeout > 0 {
		client.opts.Timeout = c.Timeout
	}
	if c.KeepAlive > 0 {
		client.dialer.KeepAlive = c.KeepAlive
		client.opts.KeepAlive = c.KeepAlive
	}
	if c.DialTimeout > 0 {
		client.dialer.Timeout = c.DialTimeout
		client.opts.Timeout = c.DialTimeout
	}
	if c.Breaker != nil {
		client.opts.Breaker = c.Breaker
		client.breaker.Reload(c.Breaker)
	}
	for uri, cfg := range c.URL {
		client.urlConf[uri] = cfg
	}
	for host, cfg := range c.Host {
		client.hostConf[host] = cfg
	}
	client.mu.Unlock()
}

// NewRequest new http request with method, uri, ip, values and headers.
func (client *Client) NewRequest(method, url string, params url.Values) (req *xhttp.Request, err error) {
	if method == xhttp.MethodGet {
		req, err = xhttp.NewRequest(xhttp.MethodGet, fmt.Sprintf("%s?%s", url, params.Encode()), nil)
	} else {
		req, err = xhttp.NewRequest(xhttp.MethodPost, url, strings.NewReader(params.Encode()))
	}
	if err != nil {
		err = pkgerr.Wrapf(err, "method:%s,uri:%s", method, url)
		return
	}
	if method == xhttp.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("User-Agent", _userAgent+" "+env.AppName)
	return
}

// Get issues a GET to the specified URL.
func (client *Client) Get(c context.Context, url string, params url.Values, res interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodGet, url, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res)
}

// Post issues a Post to the specified URL.
func (client *Client) Post(c context.Context, url string, params url.Values, res interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodPost, url, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res)
}

// RESTGet issues a RESTful GET to the specified URL.
func (client *Client) RESTGet(c context.Context, url string, params url.Values, res interface{}, v ...interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodGet, fmt.Sprintf(url, v...), params)
	if err != nil {
		return
	}
	return client.Do(c, req, res, url)
}

// RESTPost issues a RESTful Post to the specified URL.
func (client *Client) RESTPost(c context.Context, url string, params url.Values, res interface{}, v ...interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodPost, fmt.Sprintf(url, v...), params)
	if err != nil {
		return
	}
	return client.Do(c, req, res, url)
}

// Raw sends an HTTP request and returns bytes response
func (client *Client) Raw(c context.Context, req *xhttp.Request, v ...string) (bs []byte, err error) {
	var (
		ok      bool
		code    string
		cancel  func()
		resp    *xhttp.Response
		opts    *ClientOptions
		timeout time.Duration
		uri     = fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.Host, req.URL.Path)
	)
	// NOTE fix prom & config uri key.
	if len(v) == 1 {
		uri = v[0]
	}
	// breaker
	brk := client.breaker.Get(uri)
	if err = brk.Allow(); err != nil {
		code = "breaker"
		_metricClientReqCodeTotal.Inc(uri, req.Method, code)
		return
	}
	defer client.onBreaker(brk, &err)
	span, c := newClientSpan(c, req)
	defer span.Finish()
	// stat
	now := time.Now()
	defer func() {
		_metricClientReqDur.Observe(int64(time.Since(now)/time.Millisecond), uri, req.Method)
		if code != "" {
			_metricClientReqCodeTotal.Inc(uri, req.Method, code)
		}
	}()
	// get config
	// 1.url config 2.host config 3.default
	client.mu.RLock()
	if opts, ok = client.urlConf[uri]; !ok {
		if opts, ok = client.hostConf[req.Host]; !ok {
			opts = client.opts
		}
	}
	client.mu.RUnlock()
	// timeout
	deliver := true
	timeout = opts.Timeout
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := time.Until(deadline); ctimeout < timeout {
			// deliver small timeout
			timeout = ctimeout
			deliver = false
		}
	}
	if deliver {
		c, cancel = context.WithTimeout(c, timeout)
		defer cancel()
	}
	req = req.WithContext(c)
	if resp, err = client.client.Do(req); err != nil {
		err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		code = "failed"
		return
	}
	defer resp.Body.Close()
	ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))
	if resp.StatusCode >= xhttp.StatusBadRequest {
		err = pkgerr.Errorf("incorrect http status:%d host:%s, url:%s", resp.StatusCode, req.URL.Host, realURL(req))
		code = strconv.Itoa(resp.StatusCode)
		return
	}
	if bs, err = readAll(resp.Body, _minRead); err != nil {
		err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		return
	}
	return
}

// Do sends an HTTP request and returns an HTTP json response.
func (client *Client) Do(c context.Context, req *xhttp.Request, res interface{}, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil {
		if err = json.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

// JSON sends an HTTP request and returns an HTTP json response.
func (client *Client) JSON(c context.Context, req *xhttp.Request, res interface{}, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil {
		if err = json.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

// PB sends an HTTP request and returns an HTTP proto response.
func (client *Client) PB(c context.Context, req *xhttp.Request, res proto.Message, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil {
		if err = proto.Unmarshal(bs, res); err != nil {
			err = pkgerr.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

func (client *Client) onBreaker(breaker breaker.Breaker, err *error) {
	if err != nil && *err != nil {
		breaker.MarkFailed()
	} else {
		breaker.MarkSuccess()
	}
}

// realUrl return url with http://host/params.
func realURL(req *xhttp.Request) string {
	if req.Method == xhttp.MethodGet {
		return req.URL.String()
	} else if req.Method == xhttp.MethodPost {
		ru := req.URL.Path
		if req.Body != nil {
			rd, ok := req.Body.(io.Reader)
			if ok {
				buf := bytes.NewBuffer([]byte{})
				buf.ReadFrom(rd)
				ru = ru + "?" + buf.String()
			}
		}
		return ru
	}
	return req.URL.Path
}

// readAll reads from r until an error or EOF and returns the data it read
// from the internal buffer allocated with a specified capacity.
func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}
