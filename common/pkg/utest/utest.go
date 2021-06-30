/*
Package utest - contains common functions and interfaces for Unit tests
*/
package utest

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (port string) {
	var a *net.TCPAddr
	var err error
	if a, err = net.ResolveTCPAddr("tcp", ":0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return fmt.Sprintf("%d", l.Addr().(*net.TCPAddr).Port)
		}
	}
	return
}

// RoundTripper is an interface representing the ability to execute a single HTTP transaction, obtaining the Response for a given Request.
// Reference to http://hassansin.github.io/Unit-Testing-http-client-in-Go
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func NewTestClientString(s string) *http.Client {
	fn := func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       ioutil.NopCloser(bytes.NewBufferString(s)),
		}
	}
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func NewTestClientError() *http.Client {
	fn := func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       ioutil.NopCloser(errReader(0)),
		}
	}
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
